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

package gc

import (
	"context"
	"time"

	"github.com/pingcap/failpoint"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/config"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/pdutil"
	pd "github.com/tikv/pd/client"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// gcSafepointUpdateInterval is the minimum interval that CDC can update gc safepoint
var gcSafepointUpdateInterval = 1 * time.Minute

// Manager is an interface for gc manager
type Manager interface {
	// TryUpdateGCSafePoint tries to update TiCDC service GC safepoint.
	// Manager may skip update when it thinks it is too frequent.
	// Set `forceUpdate` to force Manager update.
	TryUpdateGCSafePoint(ctx context.Context, safePoint model.Ts, forceUpdate bool) error
	CheckStaleCheckpointTs(changefeedID model.ChangeFeedID, checkpointTs model.Ts) error
}

type gcManager struct {
	gcServiceID string
	pdClient    pd.Client
	pdClock     pdutil.Clock
	gcTTL       int64

	lastUpdatedTime       time.Time
	lastSucceededTime     time.Time
	lastGlobalGCSafepoint uint64
	isTiCDCBlockGC        atomic.Bool
}

// NewManager creates a new Manager.
func NewManager(gcServiceID string, pdClient pd.Client, pdClock pdutil.Clock) Manager {
	serverConfig := config.GetGlobalServerConfig()
	failpoint.Inject("InjectGcSafepointUpdateInterval", func(val failpoint.Value) {
		gcSafepointUpdateInterval = time.Duration(val.(int) * int(time.Millisecond))
	})
	return &gcManager{
		gcServiceID:       gcServiceID,
		pdClient:          pdClient,
		pdClock:           pdClock,
		lastSucceededTime: time.Now(),
		gcTTL:             serverConfig.GcTTL,
	}
}

func (m *gcManager) TryUpdateGCSafePoint(
	ctx context.Context, safePoint model.Ts, forceUpdate bool,
) error {
	if time.Since(m.lastUpdatedTime) < gcSafepointUpdateInterval && !forceUpdate {
		return nil
	}
	m.lastUpdatedTime = time.Now()

	globalGCSafepoint, err := SetServiceGCSafepoint(
		ctx, m.pdClient, m.gcServiceID, m.gcTTL, safePoint)
	if err != nil {
		log.Warn("update gc safe point failed",
			zap.String("serviceID", m.gcServiceID),
			zap.Uint64("safePointTs", safePoint),
			zap.Error(err))
		if time.Since(m.lastSucceededTime) >= time.Second*time.Duration(m.gcTTL) {
			return cerror.ErrUpdateServiceSafepointFailed.Wrap(err)
		}
		return nil
	}
	failpoint.Inject("InjectActualGCSafePoint", func(val failpoint.Value) {
		globalGCSafepoint = uint64(val.(int))
	})

	if globalGCSafepoint > safePoint {
		log.Warn("update gc safe point failed, the gc safe point is larger than checkpointTs",
			zap.Uint64("actual", globalGCSafepoint), zap.Uint64("checkpointTs", safePoint))
	} else {
		log.Info("update gc safe point success",
			zap.Uint64("globalGCSafepoint", globalGCSafepoint), zap.Uint64("gcSafePointTs", safePoint))
	}
	// if the global gc safe point is equal to the safePoint just set by TiCDC,
	// which means TiCDC is blocking the whole TiDB cluster GC.
	m.isTiCDCBlockGC.Store(globalGCSafepoint == safePoint)
	m.lastGlobalGCSafepoint = globalGCSafepoint
	m.lastSucceededTime = time.Now()
	return nil
}

func (m *gcManager) CheckStaleCheckpointTs(
	changefeedID model.ChangeFeedID, checkpointTs model.Ts,
) error {
	gcSafepointUpperBound := checkpointTs - 1
	// the data required by the TiCDC is already GCed,
	// this makes the changefeed cannot make progress anymore.
	if gcSafepointUpperBound < m.lastGlobalGCSafepoint {
		return cerror.ErrSnapshotLostByGC.GenWithStackByArgs(checkpointTs, m.lastGlobalGCSafepoint)
	}

	//if m.isTiCDCBlockGC.Load() {
	//	pdTime := m.pdClock.CurrentTime()
	//	if pdTime.Sub(
	//		oracle.GetTimeFromTS(gcSafepointUpperBound),
	//	) > time.Duration(m.gcTTL)*time.Second {
	//		return cerror.ErrGCTTLExceeded.
	//			GenWithStackByArgs(
	//				checkpointTs,
	//				changefeedID,
	//			)
	//	}
	//} else {
	//	// if `isTiCDCBlockGC` is false, it means there is another service gc
	//	// point less than the min checkpoint ts.
	//	if gcSafepointUpperBound < m.lastGlobalGCSafepoint {
	//		return cerror.ErrSnapshotLostByGC.
	//			GenWithStackByArgs(
	//				checkpointTs,
	//				m.lastGlobalGCSafepoint,
	//			)
	//	}
	//}
	//return nil
}
