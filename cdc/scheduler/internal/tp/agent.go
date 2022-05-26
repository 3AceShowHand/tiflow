// Copyright 2022 PingCAP, Inc.
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

package tp

import (
	"context"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/scheduler/internal"
	"github.com/pingcap/tiflow/cdc/scheduler/internal/tp/schedulepb"
	"github.com/pingcap/tiflow/pkg/config"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/etcd"
	"github.com/pingcap/tiflow/pkg/p2p"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

const (
	getOwnerFromEtcdTimeout         = time.Second * 5
	messageHandlerOperationsTimeout = time.Second * 5
	barrierNotAdvancingWarnDuration = time.Second * 10
	printWarnLogMinInterval         = time.Second * 1
)

var _ internal.Agent = (*agent)(nil)

type agent struct {
	trans     transport
	tableExec internal.TableExecutor

	tables       map[model.TableID]*schedulepb.TableStatus
	runningTasks map[model.TableID]*schedulepb.Message

	// maintain owner information
	etcdClient *etcd.CDCEtcdClient
	ownerInfo  *ownerInfo

	// maintain capture information
	epoch        string
	captureID    model.CaptureID
	changeFeedID model.ChangeFeedID

	// todo: these fields need to be revised
	barrierSeqs map[p2p.Topic]p2p.Seq
}

type ownerInfo struct {
	version   string
	captureID string
	revision  int64
}

func NewAgent(ctx context.Context,
	changeFeedID model.ChangeFeedID,
	messageServer *p2p.MessageServer,
	messageRouter p2p.MessageRouter,
	etcdClient *etcd.CDCEtcdClient,
	tableExecutor internal.TableExecutor) (internal.Agent, error) {
	result := &agent{
		changeFeedID: changeFeedID,
		tableExec:    tableExecutor,
		tables:       make(map[model.TableID]*schedulepb.TableStatus),
		runningTasks: make(map[model.TableID]*schedulepb.Message),
		etcdClient:   etcdClient,
	}
	trans, err := newTransport(ctx, changeFeedID, messageServer, messageRouter)
	if err != nil {
		return nil, err
	}
	result.trans = trans

	conf := config.GetGlobalServerConfig()
	flushInterval := time.Duration(conf.ProcessorFlushInterval)

	log.Debug("creating processor agent",
		zap.String("namespace", changeFeedID.Namespace),
		zap.String("changefeed", changeFeedID.ID),
		zap.Duration("sendCheckpointTsInterval", flushInterval))

	//result.Agent = base.NewBaseAgent(changeFeedID, tableExecutor, result, &base.AgentConfig{SendCheckpointTsInterval: flushInterval})

	result.etcdClient = etcdClient
	if err := result.refreshOwnerInfo(ctx); err != nil {
		return result, errors.Trace(err)
	}

	return result, nil
}

func (a *agent) refreshOwnerInfo(ctx context.Context) error {
	etcdCliCtx, cancel := context.WithTimeout(ctx, getOwnerFromEtcdTimeout)
	defer cancel()
	captureID, err := a.etcdClient.GetOwnerID(etcdCliCtx, etcd.CaptureOwnerKey)
	if err != nil {
		if err != concurrency.ErrElectionNoLeader {
			return err
		}
		// We tolerate the situation where there is no owner.
		// If we are registered in Etcd, an elected Owner will have to
		// contact us before it can schedule any table.
		log.Info("no owner found. We will wait for an owner to contact us.",
			zap.String("namespace", a.changeFeedID.Namespace),
			zap.String("changefeed", a.changeFeedID.ID),
			zap.Error(err))
		return nil
	}
	a.ownerInfo.captureID = captureID

	log.Debug("found owner",
		zap.String("namespace", a.changeFeedID.Namespace),
		zap.String("changefeed", a.changeFeedID.ID),
		zap.String("ownerID", captureID))

	revision, err := a.etcdClient.GetOwnerRevision(etcdCliCtx, captureID)
	if err != nil {
		if cerror.ErrOwnerNotFound.Equal(err) || cerror.ErrNotOwner.Equal(err) {
			// These are expected errors when no owner has been elected
			log.Info("no owner found when querying for the owner revision",
				zap.String("namespace", a.changeFeedID.Namespace),
				zap.String("changefeed", a.changeFeedID.ID),
				zap.Error(err))
			a.ownerInfo.captureID = ""
			return nil
		}
		return err
	}
	a.ownerInfo.revision = revision

	return nil
}

// Tick implement agent interface
func (a *agent) Tick(ctx context.Context) error {
	if err := a.refreshOwnerInfo(ctx); err != nil {
		return errors.Trace(err)
	}

	inboundMessages, err := a.trans.Recv(ctx)
	if err != nil {
		return errors.Trace(err)
	}
	outboundMessages, err := a.handleMessage(ctx, inboundMessages)
	if err != nil {
		return errors.Trace(err)
	}

	if err := a.trans.Send(ctx, outboundMessages); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (a *agent) handleMessage(ctx context.Context, msg []*schedulepb.Message) ([]*schedulepb.Message, error) {
	result := make([]*schedulepb.Message, 0)
	for _, message := range msg {
		header := message.GetHeader()
		if !a.updateOwnerInfo(message.GetFrom(), header.GetVersion(), header.GetOwnerRevision().Revision) {
			continue
		}

		switch message.GetMsgType() {
		case schedulepb.MsgDispatchTableRequest:
			response, err := a.handleMessageDispatchTableRequest(ctx, message.DispatchTableRequest)
			if err != nil {
				log.Warn("agent: handle dispatch table request failed", zap.Error(err))
			}
			result = append(result, response)
		case schedulepb.MsgHeartbeat:
			response, err := a.handleMessageHeartbeat()
			if err != nil {
				log.Warn("agent: handle heartbeat failed", zap.Error(err))
			}
			if response != nil {
				result = append(result, response)
			}
		case schedulepb.MsgUnknown:
		default:
			log.Warn("unknown message received")
		}
	}

	return result, nil
}

func (a *agent) handleMessageHeartbeat() (*schedulepb.Message, error) {
	// TODO: build s.tables from Heartbeat message.
	tables := make([]schedulepb.TableStatus, 0, len(a.tables))
	for _, table := range a.tables {
		// `table` should always track the latest information ?
		tables = append(tables, *table)
	}
	response := &schedulepb.HeartbeatResponse{
		Tables: tables,
		// todo (Ling Jin): how to set `IsStopping`
		IsStopping: false,
	}
	return &schedulepb.Message{
		Header:            a.newMessageHeader(),
		MsgType:           schedulepb.MsgHeartbeatResponse,
		From:              a.captureID,
		To:                a.ownerInfo.captureID,
		HeartbeatResponse: response,
	}, nil
}

func (a *agent) handleMessageDispatchTableRequest(ctx context.Context, msg *schedulepb.DispatchTableRequest) (*schedulepb.Message, error) {
	// TODO: update s.tables from DispatchTableResponse message.
	addTableRequest := msg.GetAddTable()
	if addTableRequest != nil {
		return a.handleAddTableRequest(ctx, addTableRequest)
	}

	removeTableRequest := msg.GetRemoveTable()
	if removeTableRequest != nil {
		return a.handleRemoveTableRequest(ctx, removeTableRequest)
	}

	log.Panic("agent: dispatch table request is nil")
	return nil, nil
}

//func (a *Agent) OnOwnerDispatchedTask(
//	ownerCaptureID model.CaptureID,
//	ownerRev int64,
//	tableID model.TableID,
//	startTs model.Ts,
//	isDelete bool,
//	epoch protocol.ProcessorEpoch,
//) {
//	if !a.updateOwnerInfo(ownerCaptureID, ownerRev) {
//		a.logger.Info("task from stale owner ignored",
//			zap.Int64("tableID", tableID),
//			zap.Bool("isDelete", isDelete))
//		return
//	}
//
//	a.pendingOpsMu.Lock()
//	defer a.pendingOpsMu.Unlock()
//
//	op := &agentOperation{
//		TableID:     tableID,
//		StartTs:     startTs,
//		IsDelete:    isDelete,
//		Epoch:       epoch,
//		FromOwnerID: ownerCaptureID,
//		status:      operationReceived,
//	}
//	a.pendingOps.PushBack(op)
//
//	a.logger.Info("OnOwnerDispatchedTask",
//		zap.String("ownerCaptureID", ownerCaptureID),
//		zap.Int64("ownerRev", ownerRev),
//		zap.Any("op", op))
//}

func (a *agent) handleAddTableRequest(ctx context.Context, request *schedulepb.AddTableRequest) (*schedulepb.Message, error) {
	var (
		checkpointTs model.Ts
		resolvedTs   model.Ts
	)

	tableID := request.GetTableID()
	isPrepare := !request.GetIsSecondary()
	checkpoint := request.GetCheckpoint().GetCheckpointTs()

	done, err := a.tableExec.AddTable(ctx, tableID, checkpoint, isPrepare)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if done {
		// todo: how to handle done here ?
	}

	message := &schedulepb.Message{
		Header:  a.newMessageHeader(),
		MsgType: schedulepb.MsgDispatchTableResponse,
		From:    a.captureID,
		To:      a.ownerInfo.captureID,
		DispatchTableResponse: &schedulepb.DispatchTableResponse{
			Response: &schedulepb.DispatchTableResponse_AddTable{
				AddTable: &schedulepb.AddTableResponse{
					Status: &schedulepb.TableStatus{
						TableID:    0,
						State:      0,
						Checkpoint: schedulepb.Checkpoint{},
					},
					Checkpoint: a.newCheckpoint(checkpointTs, resolvedTs, tableID),
					Reject:     false,
				},
			},
		},
	}

	return message, nil
}

func (a *agent) handleRemoveTableRequest(ctx context.Context, request *schedulepb.RemoveTableRequest) (*schedulepb.Message, error) {
	var (
		checkpointTs model.Ts
		resolvedTs   model.Ts
	)

	tableID := request.GetTableID()
	done, err := a.tableExec.RemoveTable(ctx, tableID)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if done {
		// todo: how to handle it here
	}

	message := &schedulepb.Message{
		Header:  a.newMessageHeader(),
		MsgType: schedulepb.MsgDispatchTableResponse,
		From:    a.captureID,
		To:      a.ownerInfo.captureID,
		DispatchTableResponse: &schedulepb.DispatchTableResponse{
			Response: &schedulepb.DispatchTableResponse_RemoveTable{
				RemoveTable: &schedulepb.RemoveTableResponse{
					Status: &schedulepb.TableStatus{
						TableID: tableID,
						State:   0,
						Checkpoint: schedulepb.Checkpoint{
							CheckpointTs: 0,
							ResolvedTs:   0,
							TableIDs:     []model.TableID{tableID},
						},
					},
					Checkpoint: a.newCheckpoint(checkpointTs, resolvedTs, tableID),
				},
			},
		},
	}

	return message, nil
}

// GetLastSentCheckpointTs implement agent interface
func (a *agent) GetLastSentCheckpointTs() (checkpointTs model.Ts) {
	return internal.CheckpointCannotProceed
}

// Close implement agent interface
func (a *agent) Close() error {
	return nil
}

func (a *agent) newCheckpointMessage(version string, revision int64) *schedulepb.Message {
	tableIDs := make([]model.TableID, 0, len(a.tables))
	for tableID := range a.tables {
		tableIDs = append(tableIDs, tableID)
	}
	checkpointTs, resolvedTs := a.tableExec.GetCheckpoint()
	return &schedulepb.Message{
		Header:     a.newMessageHeader(),
		MsgType:    schedulepb.MsgCheckpoint,
		From:       a.captureID,
		To:         a.ownerInfo.captureID,
		Checkpoint: a.newCheckpoint(checkpointTs, resolvedTs, tableIDs...),
	}
}

func (a *agent) newMessageHeader() *schedulepb.Message_Header {
	return &schedulepb.Message_Header{
		Version:        a.ownerInfo.version,
		OwnerRevision:  schedulepb.OwnerRevision{Revision: a.ownerInfo.revision},
		ProcessorEpoch: schedulepb.ProcessorEpoch{Epoch: a.epoch},
	}
}

func (a *agent) newCheckpoint(checkpointTs, resolvedTs model.Ts, tables ...model.TableID) *schedulepb.Checkpoint {
	return &schedulepb.Checkpoint{
		CheckpointTs: checkpointTs,
		ResolvedTs:   resolvedTs,
		TableIDs:     tables,
	}
}

// updateOwnerInfo tries to update the stored ownerInfo, and returns false if the
// owner is stale, in which case the incoming message should be ignored since
// it has come from an owner that for sure is dead.
//
// ownerCaptureID: the incoming owner's capture ID
// ownerRev: the incoming owner's revision as generated by Etcd election.
func (a *agent) updateOwnerInfo(id model.CaptureID, version string, revision int64) bool {
	// staled owner heartbeat, just ignore it.
	if a.ownerInfo.revision > revision {
		log.Info("heartbeat: from staled owner",
			zap.Any("staledOwner", ownerInfo{
				captureID: id,
				revision:  revision,
			}),
			zap.Any("owner", a.ownerInfo))
		return false
	}

	if a.ownerInfo.revision == revision && a.ownerInfo.captureID != id {
		// This panic will happen only if two messages have been received
		// with the same ownerRev but with different ownerIDs.
		// This should never happen unless the election via Etcd is buggy.
		log.Panic("owner IDs do not match",
			zap.String("expected", a.ownerInfo.captureID),
			zap.String("actual", id))
	}

	if a.ownerInfo.revision < revision {
		a.ownerInfo.captureID = id
		a.ownerInfo.revision = revision
		a.ownerInfo.version = version

		log.Info("new owner in power", zap.Any("owner", a.ownerInfo))

		// todo: shall drop all pending operations ?
		// Resets the deque so that pending operations from the previous owner
		// will not be processed.
		// Note: these pending operations have not yet been processed by the agent,
		// so it is okay to lose them.
		//a.pendingOpsMu.Lock()
		//a.pendingOps = deque.NewDeque()
		//a.pendingOpsMu.Unlock()
		//return true
	}
	return true
}
