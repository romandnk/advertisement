// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	models "github.com/romandnk/advertisement/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockAdvert is a mock of Advert interface.
type MockAdvert struct {
	ctrl     *gomock.Controller
	recorder *MockAdvertMockRecorder
}

// MockAdvertMockRecorder is the mock recorder for MockAdvert.
type MockAdvertMockRecorder struct {
	mock *MockAdvert
}

// NewMockAdvert creates a new mock instance.
func NewMockAdvert(ctrl *gomock.Controller) *MockAdvert {
	mock := &MockAdvert{ctrl: ctrl}
	mock.recorder = &MockAdvertMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdvert) EXPECT() *MockAdvertMockRecorder {
	return m.recorder
}

// CreateAdvert mocks base method.
func (m *MockAdvert) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdvert", ctx, advert)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdvert indicates an expected call of CreateAdvert.
func (mr *MockAdvertMockRecorder) CreateAdvert(ctx, advert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdvert", reflect.TypeOf((*MockAdvert)(nil).CreateAdvert), ctx, advert)
}

// DeleteAdvert mocks base method.
func (m *MockAdvert) DeleteAdvert(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdvert", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdvert indicates an expected call of DeleteAdvert.
func (mr *MockAdvertMockRecorder) DeleteAdvert(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdvert", reflect.TypeOf((*MockAdvert)(nil).DeleteAdvert), ctx, id)
}

// MockServices is a mock of Services interface.
type MockServices struct {
	ctrl     *gomock.Controller
	recorder *MockServicesMockRecorder
}

// MockServicesMockRecorder is the mock recorder for MockServices.
type MockServicesMockRecorder struct {
	mock *MockServices
}

// NewMockServices creates a new mock instance.
func NewMockServices(ctrl *gomock.Controller) *MockServices {
	mock := &MockServices{ctrl: ctrl}
	mock.recorder = &MockServicesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServices) EXPECT() *MockServicesMockRecorder {
	return m.recorder
}

// CreateAdvert mocks base method.
func (m *MockServices) CreateAdvert(ctx context.Context, advert models.Advert) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdvert", ctx, advert)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdvert indicates an expected call of CreateAdvert.
func (mr *MockServicesMockRecorder) CreateAdvert(ctx, advert interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdvert", reflect.TypeOf((*MockServices)(nil).CreateAdvert), ctx, advert)
}

// DeleteAdvert mocks base method.
func (m *MockServices) DeleteAdvert(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdvert", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdvert indicates an expected call of DeleteAdvert.
func (mr *MockServicesMockRecorder) DeleteAdvert(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdvert", reflect.TypeOf((*MockServices)(nil).DeleteAdvert), ctx, id)
}
