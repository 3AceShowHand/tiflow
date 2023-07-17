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

package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pingcap/log"
	"github.com/pingcap/tidb/br/pkg/storage"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/sink/metrics/mq"
	"github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/sink/codec/common"
	"github.com/pingcap/tiflow/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// ClaimCheck manage send message to the claim-check external storage.
type ClaimCheck struct {
	storage storage.ExternalStorage

	changefeedID model.ChangeFeedID

	// metricSendMessageDuration tracks the time duration
	// cost on send messages to the claim check external storage.
	metricSendMessageDuration prometheus.Observer
	metricSendMessageCount    prometheus.Counter
}

// NewClaimCheck return a new ClaimCheck.
func NewClaimCheck(ctx context.Context, uri string, changefeedID model.ChangeFeedID) (*ClaimCheck, error) {
	storage, err := util.GetExternalStorageFromURI(ctx, uri)
	if err != nil {
		return nil, errors.Trace(err)
	}

	log.Info("claim-check enabled",
		zap.String("namespace", changefeedID.Namespace),
		zap.String("changefeed", changefeedID.ID),
		zap.String("storageURI", uri))

	return &ClaimCheck{
		changefeedID:              changefeedID,
		storage:                   storage,
		metricSendMessageDuration: mq.ClaimCheckSendMessageDuration.WithLabelValues(changefeedID.Namespace, changefeedID.ID),
		metricSendMessageCount:    mq.ClaimCheckSendMessageCount.WithLabelValues(changefeedID.Namespace, changefeedID.ID),
	}, nil
}

// WriteMessage write message to the claim check external storage.
func (c *ClaimCheck) WriteMessage(ctx context.Context, key, value []byte, fileName string) error {
	m := common.ClaimCheckMessage{
		Key:   key,
		Value: value,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return errors.Trace(err)
	}

	start := time.Now()
	err = c.storage.WriteFile(ctx, fileName, data)
	if err != nil {
		return errors.Trace(err)
	}
	c.metricSendMessageDuration.Observe(time.Since(start).Seconds())
	c.metricSendMessageCount.Inc()
	return nil
}

// Close the claim check by clean up the metrics.
func (c *ClaimCheck) Close() {
	mq.ClaimCheckSendMessageDuration.DeleteLabelValues(c.changefeedID.Namespace, c.changefeedID.ID)
	mq.ClaimCheckSendMessageCount.DeleteLabelValues(c.changefeedID.Namespace, c.changefeedID.ID)
}