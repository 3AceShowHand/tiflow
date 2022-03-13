// Copyright 2020 PingCAP, Inc.
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

package main

import (
	"encoding/json"
	"regexp"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/etcd"
	"github.com/pingcap/tiflow/pkg/identity"
	"github.com/pingcap/tiflow/pkg/orchestrator"
	"github.com/pingcap/tiflow/pkg/orchestrator/util"
	"go.uber.org/zap"
)

type cdcReactorState struct {
	Owner              model.CaptureID
	Captures           map[model.CaptureID]*model.CaptureInfo
	ChangefeedStatuses map[identity.ChangeFeedID]*model.ChangeFeedStatus
	TaskPositions      map[identity.ChangeFeedID]map[model.CaptureID]*model.TaskPosition
	TaskStatuses       map[identity.ChangeFeedID]map[model.CaptureID]*model.TaskStatus
}

var (
	captureRegex    = regexp.MustCompile(regexp.QuoteMeta(etcd.CaptureInfoKeyPrefix) + "/(.+)")
	changefeedRegex = regexp.MustCompile(regexp.QuoteMeta(etcd.JobKeyPrefix) + "/(.+)")
	positionRegex   = regexp.MustCompile(regexp.QuoteMeta(etcd.TaskPositionKeyPrefix) + "/(.+?)/(.+)")
	statusRegex     = regexp.MustCompile(regexp.QuoteMeta(etcd.TaskStatusKeyPrefix) + "/(.+?)/(.+)")
)

func newCDCReactorState() *cdcReactorState {
	return &cdcReactorState{
		Captures:           make(map[model.CaptureID]*model.CaptureInfo),
		ChangefeedStatuses: make(map[identity.ChangeFeedID]*model.ChangeFeedStatus),
		TaskPositions:      make(map[identity.ChangeFeedID]map[model.CaptureID]*model.TaskPosition),
		TaskStatuses:       make(map[identity.ChangeFeedID]map[model.CaptureID]*model.TaskStatus),
	}
}

func (s *cdcReactorState) Update(key util.EtcdKey, value []byte, isInit bool) error {
	if key.String() == etcd.CaptureOwnerKey {
		if value == nil {
			log.Info("Owner lost", zap.String("oldOwner", s.Owner))
			return nil
		}

		log.Info("Owner updated", zap.String("oldOwner", s.Owner),
			zap.ByteString("newOwner", value))
		s.Owner = string(value)
		return nil
	}

	if matches := captureRegex.FindSubmatch(key.Bytes()); matches != nil {
		captureID := string(matches[1])

		if value == nil {
			log.Info("Capture deleted",
				zap.String("captureID", captureID),
				zap.Reflect("oldCapture", s.Captures[captureID]))

			delete(s.Captures, captureID)
			return nil
		}

		var newCaptureInfo model.CaptureInfo
		err := json.Unmarshal(value, &newCaptureInfo)
		if err != nil {
			return errors.Trace(err)
		}

		if oldCaptureInfo, ok := s.Captures[captureID]; ok {
			log.Info("Capture updated",
				zap.String("captureID", captureID),
				zap.Reflect("oldCapture", oldCaptureInfo),
				zap.Reflect("newCapture", newCaptureInfo))
		} else {
			log.Info("Capture added",
				zap.String("captureID", captureID),
				zap.Reflect("newCapture", newCaptureInfo))
		}

		s.Captures[captureID] = &newCaptureInfo
		return nil
	}

	if matches := changefeedRegex.FindSubmatch(key.Bytes()); matches != nil {
		changefeedID := string(matches[1])

		if value == nil {
			log.Info("Changefeed deleted",
				zap.String("changefeed", changefeedID),
				zap.Reflect("oldChangefeed", s.ChangefeedStatuses))

			delete(s.ChangefeedStatuses, changefeedID)
			return nil
		}

		var newChangefeedStatus model.ChangeFeedStatus
		err := json.Unmarshal(value, &newChangefeedStatus)
		if err != nil {
			return errors.Trace(err)
		}

		if oldChangefeedInfo, ok := s.ChangefeedStatuses[changefeedID]; ok {
			log.Info("Changefeed updated",
				zap.String("changefeed", changefeedID),
				zap.Reflect("oldChangefeed", oldChangefeedInfo),
				zap.Reflect("newChangefeed", newChangefeedStatus))
		} else {
			log.Info("Changefeed added",
				zap.String("changefeed", changefeedID),
				zap.Reflect("newChangefeed", newChangefeedStatus))
		}

		s.ChangefeedStatuses[changefeedID] = &newChangefeedStatus

		return nil
	}

	if matches := positionRegex.FindSubmatch(key.Bytes()); matches != nil {
		captureID := string(matches[1])
		changefeedID := string(matches[2])

		if value == nil {
			log.Info("Position deleted",
				zap.String("captureID", captureID),
				zap.String("changefeed", changefeedID),
				zap.Reflect("oldPosition", s.TaskPositions[changefeedID][captureID]))

			delete(s.TaskPositions[changefeedID], captureID)
			if len(s.TaskPositions[changefeedID]) == 0 {
				delete(s.TaskPositions, changefeedID)
			}

			return nil
		}

		var newTaskPosition model.TaskPosition
		err := json.Unmarshal(value, &newTaskPosition)
		if err != nil {
			return errors.Trace(err)
		}

		if _, ok := s.TaskPositions[changefeedID]; !ok {
			s.TaskPositions[changefeedID] = make(map[model.CaptureID]*model.TaskPosition)
		}

		if position, ok := s.TaskPositions[changefeedID][captureID]; ok {
			log.Info("Position updated",
				zap.String("captureID", captureID),
				zap.String("changefeed", changefeedID),
				zap.Reflect("oldPosition", position),
				zap.Reflect("newPosition", newTaskPosition))
		} else {
			log.Info("Position created",
				zap.String("captureID", captureID),
				zap.String("changefeed", changefeedID),
				zap.Reflect("newPosition", newTaskPosition))
		}

		s.TaskPositions[changefeedID][captureID] = &newTaskPosition

		return nil
	}

	if matches := statusRegex.FindSubmatch(key.Bytes()); matches != nil {
		captureID := string(matches[1])
		changefeedID := string(matches[2])

		if value == nil {
			log.Info("Status deleted",
				zap.String("captureID", captureID),
				zap.String("changefeed", changefeedID),
				zap.Reflect("oldStatus", s.TaskStatuses[changefeedID][captureID]))

			delete(s.TaskStatuses[changefeedID], captureID)
			if len(s.TaskStatuses[changefeedID]) == 0 {
				delete(s.TaskStatuses, changefeedID)
			}

			return nil
		}

		var newTaskStatus model.TaskStatus
		err := json.Unmarshal(value, &newTaskStatus)
		if err != nil {
			return errors.Trace(err)
		}

		if _, ok := s.TaskStatuses[changefeedID]; !ok {
			s.TaskStatuses[changefeedID] = make(map[model.CaptureID]*model.TaskStatus)
		}

		if status, ok := s.TaskStatuses[changefeedID][captureID]; ok {
			log.Info("Status updated",
				zap.String("captureID", captureID),
				zap.String("changefeed", changefeedID),
				zap.Reflect("oldStatus", status),
				zap.Reflect("newStatus", newTaskStatus))
		} else {
			log.Info("Status updated",
				zap.String("captureID", captureID),
				zap.String("changefeed", changefeedID),
				zap.Reflect("newStatus", newTaskStatus))
		}

		s.TaskStatuses[changefeedID][captureID] = &newTaskStatus

		return nil
	}

	log.Debug("Etcd operation ignored", zap.String("key", key.String()), zap.ByteString("value", value))
	return nil
}

func (s *cdcReactorState) GetPatches() [][]orchestrator.DataPatch {
	return nil
}
