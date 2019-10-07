// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpack/pack/commands (interfaces: PackClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	lifecycle "github.com/buildpack/lifecycle"
	gomock "github.com/golang/mock/gomock"

	pack "github.com/buildpack/pack"
)

// MockPackClient is a mock of PackClient interface
type MockPackClient struct {
	ctrl     *gomock.Controller
	recorder *MockPackClientMockRecorder
}

// MockPackClientMockRecorder is the mock recorder for MockPackClient
type MockPackClientMockRecorder struct {
	mock *MockPackClient
}

// NewMockPackClient creates a new mock instance
func NewMockPackClient(ctrl *gomock.Controller) *MockPackClient {
	mock := &MockPackClient{ctrl: ctrl}
	mock.recorder = &MockPackClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPackClient) EXPECT() *MockPackClientMockRecorder {
	return m.recorder
}

// Build mocks base method
func (m *MockPackClient) Build(arg0 context.Context, arg1 pack.BuildOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Build indicates an expected call of Build
func (mr *MockPackClientMockRecorder) Build(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockPackClient)(nil).Build), arg0, arg1)
}

// CreateBuilder mocks base method
func (m *MockPackClient) CreateBuilder(arg0 context.Context, arg1 pack.CreateBuilderOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBuilder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBuilder indicates an expected call of CreateBuilder
func (mr *MockPackClientMockRecorder) CreateBuilder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBuilder", reflect.TypeOf((*MockPackClient)(nil).CreateBuilder), arg0, arg1)
}

// InspectBuilder mocks base method
func (m *MockPackClient) InspectBuilder(arg0 string, arg1 bool) (*pack.BuilderInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InspectBuilder", arg0, arg1)
	ret0, _ := ret[0].(*pack.BuilderInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InspectBuilder indicates an expected call of InspectBuilder
func (mr *MockPackClientMockRecorder) InspectBuilder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InspectBuilder", reflect.TypeOf((*MockPackClient)(nil).InspectBuilder), arg0, arg1)
}

// InspectImage mocks base method
func (m *MockPackClient) InspectImage(arg0 string, arg1 bool) (*pack.ImageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InspectImage", arg0, arg1)
	ret0, _ := ret[0].(*pack.ImageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InspectImage indicates an expected call of InspectImage
func (mr *MockPackClientMockRecorder) InspectImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InspectImage", reflect.TypeOf((*MockPackClient)(nil).InspectImage), arg0, arg1)
}

// Rebase mocks base method
func (m *MockPackClient) Rebase(arg0 context.Context, arg1 lifecycle.Rebaser, arg2 pack.RebaseOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rebase", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rebase indicates an expected call of Rebase
func (mr *MockPackClientMockRecorder) Rebase(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebase", reflect.TypeOf((*MockPackClient)(nil).Rebase), arg0, arg1, arg2)
}
