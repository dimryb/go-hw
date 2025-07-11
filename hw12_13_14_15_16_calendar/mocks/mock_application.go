// Code generated by MockGen. DO NOT EDIT.
// Source: application.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	types "github.com/dimryb/go-hw/hw12_13_14_15_calendar/internal/types"
	gomock "github.com/golang/mock/gomock"
)

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// CreateEvent mocks base method.
func (m *MockApplication) CreateEvent(arg0 context.Context, arg1 types.Event) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEvent", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEvent indicates an expected call of CreateEvent.
func (mr *MockApplicationMockRecorder) CreateEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEvent", reflect.TypeOf((*MockApplication)(nil).CreateEvent), arg0, arg1)
}

// DeleteEvent mocks base method.
func (m *MockApplication) DeleteEvent(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEvent indicates an expected call of DeleteEvent.
func (mr *MockApplicationMockRecorder) DeleteEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEvent", reflect.TypeOf((*MockApplication)(nil).DeleteEvent), arg0, arg1)
}

// DeleteOlderThan mocks base method.
func (m *MockApplication) DeleteOlderThan(arg0 context.Context, arg1 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOlderThan", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOlderThan indicates an expected call of DeleteOlderThan.
func (mr *MockApplicationMockRecorder) DeleteOlderThan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOlderThan", reflect.TypeOf((*MockApplication)(nil).DeleteOlderThan), arg0, arg1)
}

// GetEventByID mocks base method.
func (m *MockApplication) GetEventByID(arg0 context.Context, arg1 string) (types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventByID", arg0, arg1)
	ret0, _ := ret[0].(types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventByID indicates an expected call of GetEventByID.
func (mr *MockApplicationMockRecorder) GetEventByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventByID", reflect.TypeOf((*MockApplication)(nil).GetEventByID), arg0, arg1)
}

// ListEvents mocks base method.
func (m *MockApplication) ListEvents(arg0 context.Context) ([]types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEvents", arg0)
	ret0, _ := ret[0].([]types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEvents indicates an expected call of ListEvents.
func (mr *MockApplicationMockRecorder) ListEvents(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEvents", reflect.TypeOf((*MockApplication)(nil).ListEvents), arg0)
}

// ListEventsByUser mocks base method.
func (m *MockApplication) ListEventsByUser(arg0 context.Context, arg1 string) ([]types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsByUser", arg0, arg1)
	ret0, _ := ret[0].([]types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsByUser indicates an expected call of ListEventsByUser.
func (mr *MockApplicationMockRecorder) ListEventsByUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsByUser", reflect.TypeOf((*MockApplication)(nil).ListEventsByUser), arg0, arg1)
}

// ListEventsByUserInRange mocks base method.
func (m *MockApplication) ListEventsByUserInRange(arg0 context.Context, arg1 string, arg2, arg3 time.Time) ([]types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsByUserInRange", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsByUserInRange indicates an expected call of ListEventsByUserInRange.
func (mr *MockApplicationMockRecorder) ListEventsByUserInRange(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsByUserInRange", reflect.TypeOf((*MockApplication)(nil).ListEventsByUserInRange), arg0, arg1, arg2, arg3)
}

// ListEventsDueBefore mocks base method.
func (m *MockApplication) ListEventsDueBefore(arg0 context.Context, arg1 time.Time) ([]types.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEventsDueBefore", arg0, arg1)
	ret0, _ := ret[0].([]types.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEventsDueBefore indicates an expected call of ListEventsDueBefore.
func (mr *MockApplicationMockRecorder) ListEventsDueBefore(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEventsDueBefore", reflect.TypeOf((*MockApplication)(nil).ListEventsDueBefore), arg0, arg1)
}

// UpdateEvent mocks base method.
func (m *MockApplication) UpdateEvent(arg0 context.Context, arg1 types.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEvent indicates an expected call of UpdateEvent.
func (mr *MockApplicationMockRecorder) UpdateEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEvent", reflect.TypeOf((*MockApplication)(nil).UpdateEvent), arg0, arg1)
}
