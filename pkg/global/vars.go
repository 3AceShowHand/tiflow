package global

import (
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/sorter/db/system"
	"github.com/pingcap/tiflow/pkg/etcd"
	"github.com/pingcap/tiflow/pkg/p2p"
)

// GlobalVars contains some vars which can be used anywhere in a pipeline
// the lifecycle of vars in the GlobalVars should be aligned with the ticdc server process.
// All field in Vars should be READ-ONLY and THREAD-SAFE
type Vars struct {
	CaptureInfo      *model.CaptureInfo
	EtcdClient       etcd.CDCEtcdClient
	TableActorSystem *system.System
	SorterSystem     *ssystem.System

	// OwnerRevision is the Etcd revision when the owner got elected.
	OwnerRevision int64

	// MessageServer and MessageRouter are for peer-messaging
	MessageServer *p2p.MessageServer
	MessageRouter p2p.MessageRouter
}

func new() *Vars {
	return &Vars{}
}
