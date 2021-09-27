// Copyright 2021 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package processor

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/failpoint"
	"github.com/pingcap/log"
	"github.com/pingcap/ticdc/cdc/model"
	tablepipeline "github.com/pingcap/ticdc/cdc/processor/pipeline"
	cdcContext "github.com/pingcap/ticdc/pkg/context"
	cerrors "github.com/pingcap/ticdc/pkg/errors"
	"github.com/pingcap/ticdc/pkg/orchestrator"
	"go.uber.org/zap"
)

type commandType int

const (
	commandTypeUnknown commandType = iota //nolint:varcheck,deadcode
	commandTypeClose
	commandTypeWriteDebugInfo
)

type command struct {
	tp      commandType
	payload io.Writer
	done    chan struct{}
}

// Manager is a manager of processor, which maintains the state and behavior of processors
type Manager struct {
	processors map[model.ChangeFeedID]*processor

	commandQueue chan *command

	newProcessor func(cdcContext.Context) *processor
}

// NewManager creates a new processor manager
func NewManager() *Manager {
	return &Manager{
		processors:   make(map[model.ChangeFeedID]*processor),
		commandQueue: make(chan *command, 4),
		newProcessor: newProcessor,
	}
}

// NewManager4Test creates a new processor manager for test
func NewManager4Test(
	createTablePipeline func(ctx cdcContext.Context, tableID model.TableID, replicaInfo *model.TableReplicaInfo) (tablepipeline.TablePipeline, error),
) *Manager {
	m := NewManager()
	m.newProcessor = func(ctx cdcContext.Context) *processor {
		return newProcessor4Test(ctx, createTablePipeline)
	}
	return m
}

// Tick implements the `orchestrator.State` interface
// the `state` parameter is sent by the etcd worker, the `state` must be a snapshot of KVs in etcd
// the Tick function of Manager create or remove processor instances according to the specified `state`, or pass the `state` to processor instances
func (m *Manager) Tick(stdCtx context.Context, state orchestrator.ReactorState) (nextState orchestrator.ReactorState, err error) {
	ctx := stdCtx.(cdcContext.Context)
	globalState := state.(*orchestrator.GlobalReactorState)
	if err := m.handleCommand(); err != nil {
		return state, err
	}
	captureID := ctx.GlobalVars().CaptureInfo.ID
	var inactiveChangefeedCount int
	for changefeedID, changefeedState := range globalState.Changefeeds {
		if !changefeedState.Active(captureID) {
			inactiveChangefeedCount++
			m.closeProcessor(changefeedID)
			continue
		}
		ctx := cdcContext.WithChangefeedVars(ctx, &cdcContext.ChangefeedVars{
			ID:   changefeedID,
			Info: changefeedState.Info,
		})
		processor, exist := m.processors[changefeedID]
		if !exist {
			if changefeedState.Status.AdminJobType.IsStopState() {
				continue
			}
			taskStatus := changefeedState.TaskStatuses[captureID]
			if taskStatus == nil || taskStatus.AdminJobType.IsStopState() {
				continue
			}
			// the processor should start after at least one table has been added to this capture
			if len(taskStatus.Tables) == 0 && len(taskStatus.Operation) == 0 {
				continue
			}
			failpoint.Inject("processorManagerHandleNewChangefeedDelay", nil)
			processor = m.newProcessor(ctx)
			m.processors[changefeedID] = processor
		}
		if _, err := processor.Tick(ctx, changefeedState); err != nil {
			m.closeProcessor(changefeedID)
			if cerrors.ErrReactorFinished.Equal(errors.Cause(err)) {
				continue
			}
			return state, errors.Trace(err)
		}
	}
	// check if the processors in memory is leaked
	m.evictProcessor(globalState, inactiveChangefeedCount)

	return state, nil
}

// evict Processor maintained in memory, if it's not in etcd anymore.
func (m *Manager) evictProcessor(state *orchestrator.GlobalReactorState, inActive int) {
	if len(m.processors)+inActive == len(state.Changefeeds) {
		return
	}

	for changefeedID := range m.processors {
		if _, exist := state.Changefeeds[changefeedID]; !exist {
			m.closeProcessor(changefeedID)
		}
	}
}

func (m *Manager) closeProcessor(changefeedID model.ChangeFeedID) {
	if processor, exist := m.processors[changefeedID]; exist {
		err := processor.Close()
		if err != nil {
			log.Warn("failed to close processor", zap.Error(err))
		}
		delete(m.processors, changefeedID)
	}
}

// AsyncClose sends a close signal to Manager and closing all processors
func (m *Manager) AsyncClose() {
	m.sendCommand(commandTypeClose, nil)
}

// WriteDebugInfo write the debug info to Writer
func (m *Manager) WriteDebugInfo(w io.Writer) {
	timeout := time.Second * 3
	done := m.sendCommand(commandTypeWriteDebugInfo, w)
	// wait the debug info printed
	select {
	case <-done:
	case <-time.After(timeout):
		fmt.Fprintf(w, "failed to print debug info for processor\n")
	}
}

func (m *Manager) sendCommand(tp commandType, payload io.Writer) chan struct{} {
	timeout := time.Second * 3
	cmd := &command{tp: tp, payload: payload, done: make(chan struct{})}
	select {
	case m.commandQueue <- cmd:
	case <-time.After(timeout):
		close(cmd.done)
		log.Warn("the command queue is full, ignore this command", zap.Any("command", cmd))
	}
	return cmd.done
}

func (m *Manager) handleCommand() error {
	var cmd *command
	select {
	case cmd = <-m.commandQueue:
	default:
		return nil
	}
	defer close(cmd.done)
	switch cmd.tp {
	case commandTypeClose:
		for changefeedID := range m.processors {
			m.closeProcessor(changefeedID)
		}
		return cerrors.ErrReactorFinished
	case commandTypeWriteDebugInfo:
		m.writeDebugInfo(cmd.payload)
	default:
		log.Warn("Unknown command in processor manager", zap.Any("command", cmd))
	}
	return nil
}

func (m *Manager) writeDebugInfo(w io.Writer) {
	for changefeedID, processor := range m.processors {
		fmt.Fprintf(w, "changefeedID: %s\n", changefeedID)
		processor.WriteDebugInfo(w)
		fmt.Fprintf(w, "\n")
	}
}
