// Code generated by MockGen. DO NOT EDIT.
// Source: result.go
//
// Generated by this command:
//
//	mockgen -source=result.go -destination=/Users/gda/Documents/github/qase-tms/qasectl/internal/service/result/mocks/result.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	result "github.com/qase-tms/qasectl/internal/models/result"
	gomock "go.uber.org/mock/gomock"
)

// Mockclient is a mock of client interface.
type Mockclient struct {
	ctrl     *gomock.Controller
	recorder *MockclientMockRecorder
}

// MockclientMockRecorder is the mock recorder for Mockclient.
type MockclientMockRecorder struct {
	mock *Mockclient
}

// NewMockclient creates a new mock instance.
func NewMockclient(ctrl *gomock.Controller) *Mockclient {
	mock := &Mockclient{ctrl: ctrl}
	mock.recorder = &MockclientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockclient) EXPECT() *MockclientMockRecorder {
	return m.recorder
}

// UploadData mocks base method.
func (m *Mockclient) UploadData(ctx context.Context, project string, runID int64, results []result.Result) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadData", ctx, project, runID, results)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadData indicates an expected call of UploadData.
func (mr *MockclientMockRecorder) UploadData(ctx, project, runID, results any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadData", reflect.TypeOf((*Mockclient)(nil).UploadData), ctx, project, runID, results)
}

// MockParser is a mock of Parser interface.
type MockParser struct {
	ctrl     *gomock.Controller
	recorder *MockParserMockRecorder
}

// MockParserMockRecorder is the mock recorder for MockParser.
type MockParserMockRecorder struct {
	mock *MockParser
}

// NewMockParser creates a new mock instance.
func NewMockParser(ctrl *gomock.Controller) *MockParser {
	mock := &MockParser{ctrl: ctrl}
	mock.recorder = &MockParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParser) EXPECT() *MockParserMockRecorder {
	return m.recorder
}

// Parse mocks base method.
func (m *MockParser) Parse() ([]result.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse")
	ret0, _ := ret[0].([]result.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse.
func (mr *MockParserMockRecorder) Parse() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockParser)(nil).Parse))
}

// MockrunService is a mock of runService interface.
type MockrunService struct {
	ctrl     *gomock.Controller
	recorder *MockrunServiceMockRecorder
}

// MockrunServiceMockRecorder is the mock recorder for MockrunService.
type MockrunServiceMockRecorder struct {
	mock *MockrunService
}

// NewMockrunService creates a new mock instance.
func NewMockrunService(ctrl *gomock.Controller) *MockrunService {
	mock := &MockrunService{ctrl: ctrl}
	mock.recorder = &MockrunServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockrunService) EXPECT() *MockrunServiceMockRecorder {
	return m.recorder
}

// CompleteRun mocks base method.
func (m *MockrunService) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteRun", ctx, projectCode, runId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteRun indicates an expected call of CompleteRun.
func (mr *MockrunServiceMockRecorder) CompleteRun(ctx, projectCode, runId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteRun", reflect.TypeOf((*MockrunService)(nil).CompleteRun), ctx, projectCode, runId)
}

// CreateRun mocks base method.
func (m_2 *MockrunService) CreateRun(ctx context.Context, p, t, d, e string, m, plan int64, tags []string) (int64, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "CreateRun", ctx, p, t, d, e, m, plan, tags)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRun indicates an expected call of CreateRun.
func (mr *MockrunServiceMockRecorder) CreateRun(ctx, p, t, d, e, m, plan, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRun", reflect.TypeOf((*MockrunService)(nil).CreateRun), ctx, p, t, d, e, m, plan, tags)
}
