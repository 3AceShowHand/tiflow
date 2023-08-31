// Code generated by MockGen. DO NOT EDIT.
// Source: cdc/owner/status_provider.go

// Package mock_owner is a generated GoMock package.
package mock_owner

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/pingcap/tiflow/cdc/model"
)

// MockStatusProvider is a mock of StatusProvider interface.
type MockStatusProvider struct {
	ctrl     *gomock.Controller
	recorder *MockStatusProviderMockRecorder
}

// MockStatusProviderMockRecorder is the mock recorder for MockStatusProvider.
type MockStatusProviderMockRecorder struct {
	mock *MockStatusProvider
}

// NewMockStatusProvider creates a new mock instance.
func NewMockStatusProvider(ctrl *gomock.Controller) *MockStatusProvider {
	mock := &MockStatusProvider{ctrl: ctrl}
	mock.recorder = &MockStatusProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStatusProvider) EXPECT() *MockStatusProviderMockRecorder {
	return m.recorder
}

// GetAllChangeFeedInfo mocks base method.
func (m *MockStatusProvider) GetAllChangeFeedInfo(ctx context.Context) (map[model.ChangeFeedID]*model.ChangeFeedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChangeFeedInfo", ctx)
	ret0, _ := ret[0].(map[model.ChangeFeedID]*model.ChangeFeedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllChangeFeedInfo indicates an expected call of GetAllChangeFeedInfo.
func (mr *MockStatusProviderMockRecorder) GetAllChangeFeedInfo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChangeFeedInfo", reflect.TypeOf((*MockStatusProvider)(nil).GetAllChangeFeedInfo), ctx)
}

// GetAllChangeFeedStatuses mocks base method.
func (m *MockStatusProvider) GetAllChangeFeedStatuses(ctx context.Context) (map[model.ChangeFeedID]*model.ChangeFeedStatusForAPI, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllChangeFeedStatuses", ctx)
	ret0, _ := ret[0].(map[model.ChangeFeedID]*model.ChangeFeedStatusForAPI)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllChangeFeedStatuses indicates an expected call of GetAllChangeFeedStatuses.
func (mr *MockStatusProviderMockRecorder) GetAllChangeFeedStatuses(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllChangeFeedStatuses", reflect.TypeOf((*MockStatusProvider)(nil).GetAllChangeFeedStatuses), ctx)
}

// GetAllTaskStatuses mocks base method.
func (m *MockStatusProvider) GetAllTaskStatuses(ctx context.Context, changefeedID model.ChangeFeedID) (map[model.CaptureID]*model.TaskStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTaskStatuses", ctx, changefeedID)
	ret0, _ := ret[0].(map[model.CaptureID]*model.TaskStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTaskStatuses indicates an expected call of GetAllTaskStatuses.
func (mr *MockStatusProviderMockRecorder) GetAllTaskStatuses(ctx, changefeedID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTaskStatuses", reflect.TypeOf((*MockStatusProvider)(nil).GetAllTaskStatuses), ctx, changefeedID)
}

// GetCaptures mocks base method.
func (m *MockStatusProvider) GetCaptures(ctx context.Context) ([]*model.CaptureInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCaptures", ctx)
	ret0, _ := ret[0].([]*model.CaptureInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCaptures indicates an expected call of GetCaptures.
func (mr *MockStatusProviderMockRecorder) GetCaptures(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCaptures", reflect.TypeOf((*MockStatusProvider)(nil).GetCaptures), ctx)
}

// GetChangeFeedInfo mocks base method.
func (m *MockStatusProvider) GetChangeFeedInfo(ctx context.Context, changefeedID model.ChangeFeedID) (*model.ChangeFeedInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChangeFeedInfo", ctx, changefeedID)
	ret0, _ := ret[0].(*model.ChangeFeedInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChangeFeedInfo indicates an expected call of GetChangeFeedInfo.
func (mr *MockStatusProviderMockRecorder) GetChangeFeedInfo(ctx, changefeedID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChangeFeedInfo", reflect.TypeOf((*MockStatusProvider)(nil).GetChangeFeedInfo), ctx, changefeedID)
}

// GetChangeFeedStatus mocks base method.
func (m *MockStatusProvider) GetChangeFeedStatus(ctx context.Context, changefeedID model.ChangeFeedID) (*model.ChangeFeedStatusForAPI, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChangeFeedStatus", ctx, changefeedID)
	ret0, _ := ret[0].(*model.ChangeFeedStatusForAPI)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChangeFeedStatus indicates an expected call of GetChangeFeedStatus.
func (mr *MockStatusProviderMockRecorder) GetChangeFeedStatus(ctx, changefeedID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChangeFeedStatus", reflect.TypeOf((*MockStatusProvider)(nil).GetChangeFeedStatus), ctx, changefeedID)
}

// GetProcessors mocks base method.
func (m *MockStatusProvider) GetProcessors(ctx context.Context) ([]*model.ProcInfoSnap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProcessors", ctx)
	ret0, _ := ret[0].([]*model.ProcInfoSnap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProcessors indicates an expected call of GetProcessors.
func (mr *MockStatusProviderMockRecorder) GetProcessors(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProcessors", reflect.TypeOf((*MockStatusProvider)(nil).GetProcessors), ctx)
}

// IsChangefeedOwner mocks base method.
func (m *MockStatusProvider) IsChangefeedOwner(ctx context.Context, id model.ChangeFeedID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsChangefeedOwner", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsChangefeedOwner indicates an expected call of IsChangefeedOwner.
func (mr *MockStatusProviderMockRecorder) IsChangefeedOwner(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsChangefeedOwner", reflect.TypeOf((*MockStatusProvider)(nil).IsChangefeedOwner), ctx, id)
}

// IsHealthy mocks base method.
func (m *MockStatusProvider) IsHealthy(ctx context.Context) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsHealthy", ctx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsHealthy indicates an expected call of IsHealthy.
func (mr *MockStatusProviderMockRecorder) IsHealthy(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsHealthy", reflect.TypeOf((*MockStatusProvider)(nil).IsHealthy), ctx)
}
