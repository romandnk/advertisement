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

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockUser) SignIn(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockUserMockRecorder) SignIn(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockUser)(nil).SignIn), ctx, email, password)
}

// SignUp mocks base method.
func (m *MockUser) SignUp(ctx context.Context, user models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockUserMockRecorder) SignUp(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockUser)(nil).SignUp), ctx, user)
}

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

// GetAdvertByID mocks base method.
func (m *MockAdvert) GetAdvertByID(ctx context.Context, id string) (models.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertByID", ctx, id)
	ret0, _ := ret[0].(models.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertByID indicates an expected call of GetAdvertByID.
func (mr *MockAdvertMockRecorder) GetAdvertByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertByID", reflect.TypeOf((*MockAdvert)(nil).GetAdvertByID), ctx, id)
}

// MockImage is a mock of Image interface.
type MockImage struct {
	ctrl     *gomock.Controller
	recorder *MockImageMockRecorder
}

// MockImageMockRecorder is the mock recorder for MockImage.
type MockImageMockRecorder struct {
	mock *MockImage
}

// NewMockImage creates a new mock instance.
func NewMockImage(ctrl *gomock.Controller) *MockImage {
	mock := &MockImage{ctrl: ctrl}
	mock.recorder = &MockImageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImage) EXPECT() *MockImageMockRecorder {
	return m.recorder
}

// GetImageByID mocks base method.
func (m *MockImage) GetImageByID(ctx context.Context, id string) (models.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageByID", ctx, id)
	ret0, _ := ret[0].(models.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageByID indicates an expected call of GetImageByID.
func (mr *MockImageMockRecorder) GetImageByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageByID", reflect.TypeOf((*MockImage)(nil).GetImageByID), ctx, id)
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

// GetAdvertByID mocks base method.
func (m *MockServices) GetAdvertByID(ctx context.Context, id string) (models.Advert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvertByID", ctx, id)
	ret0, _ := ret[0].(models.Advert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvertByID indicates an expected call of GetAdvertByID.
func (mr *MockServicesMockRecorder) GetAdvertByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvertByID", reflect.TypeOf((*MockServices)(nil).GetAdvertByID), ctx, id)
}

// GetImageByID mocks base method.
func (m *MockServices) GetImageByID(ctx context.Context, id string) (models.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageByID", ctx, id)
	ret0, _ := ret[0].(models.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageByID indicates an expected call of GetImageByID.
func (mr *MockServicesMockRecorder) GetImageByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageByID", reflect.TypeOf((*MockServices)(nil).GetImageByID), ctx, id)
}

// SignIn mocks base method.
func (m *MockServices) SignIn(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockServicesMockRecorder) SignIn(ctx, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockServices)(nil).SignIn), ctx, email, password)
}

// SignUp mocks base method.
func (m *MockServices) SignUp(ctx context.Context, user models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", ctx, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockServicesMockRecorder) SignUp(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockServices)(nil).SignUp), ctx, user)
}
