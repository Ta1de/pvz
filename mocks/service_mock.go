package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	response "pvz/internal/api/response"
	model "pvz/internal/repository/model"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
	isgomock struct{}
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

// CreateUser mocks base method.
func (m *MockUser) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserMockRecorder) CreateUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUser)(nil).CreateUser), ctx, user)
}

// DummyLogin mocks base method.
func (m *MockUser) DummyLogin(ctx context.Context, role string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DummyLogin", ctx, role)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DummyLogin indicates an expected call of DummyLogin.
func (mr *MockUserMockRecorder) DummyLogin(ctx, role any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DummyLogin", reflect.TypeOf((*MockUser)(nil).DummyLogin), ctx, role)
}

// LoginUser mocks base method.
func (m *MockUser) LoginUser(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserMockRecorder) LoginUser(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUser)(nil).LoginUser), ctx, email, password)
}

// MockPvz is a mock of Pvz interface.
type MockPvz struct {
	ctrl     *gomock.Controller
	recorder *MockPvzMockRecorder
	isgomock struct{}
}

// MockPvzMockRecorder is the mock recorder for MockPvz.
type MockPvzMockRecorder struct {
	mock *MockPvz
}

// NewMockPvz creates a new mock instance.
func NewMockPvz(ctrl *gomock.Controller) *MockPvz {
	mock := &MockPvz{ctrl: ctrl}
	mock.recorder = &MockPvzMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPvz) EXPECT() *MockPvzMockRecorder {
	return m.recorder
}

// CreatePvz mocks base method.
func (m *MockPvz) CreatePvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePvz", ctx, pvz)
	ret0, _ := ret[0].(model.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePvz indicates an expected call of CreatePvz.
func (mr *MockPvzMockRecorder) CreatePvz(ctx, pvz any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePvz", reflect.TypeOf((*MockPvz)(nil).CreatePvz), ctx, pvz)
}

// GetPvzList mocks base method.
func (m *MockPvz) GetPvzList(ctx context.Context, limit, offset int, startDate, endDate *time.Time) ([]response.PvzFullResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPvzList", ctx, limit, offset, startDate, endDate)
	ret0, _ := ret[0].([]response.PvzFullResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPvzList indicates an expected call of GetPvzList.
func (mr *MockPvzMockRecorder) GetPvzList(ctx, limit, offset, startDate, endDate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPvzList", reflect.TypeOf((*MockPvz)(nil).GetPvzList), ctx, limit, offset, startDate, endDate)
}

// MockReception is a mock of Reception interface.
type MockReception struct {
	ctrl     *gomock.Controller
	recorder *MockReceptionMockRecorder
	isgomock struct{}
}

// MockReceptionMockRecorder is the mock recorder for MockReception.
type MockReceptionMockRecorder struct {
	mock *MockReception
}

// NewMockReception creates a new mock instance.
func NewMockReception(ctrl *gomock.Controller) *MockReception {
	mock := &MockReception{ctrl: ctrl}
	mock.recorder = &MockReceptionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReception) EXPECT() *MockReceptionMockRecorder {
	return m.recorder
}

// CloseReception mocks base method.
func (m *MockReception) CloseReception(ctx context.Context, pvzId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseReception", ctx, pvzId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseReception indicates an expected call of CloseReception.
func (mr *MockReceptionMockRecorder) CloseReception(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseReception", reflect.TypeOf((*MockReception)(nil).CloseReception), ctx, pvzId)
}

// CreateReception mocks base method.
func (m *MockReception) CreateReception(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReception", ctx, pvzId)
	ret0, _ := ret[0].(model.Reception)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReception indicates an expected call of CreateReception.
func (mr *MockReceptionMockRecorder) CreateReception(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReception", reflect.TypeOf((*MockReception)(nil).CreateReception), ctx, pvzId)
}

// MockProduct is a mock of Product interface.
type MockProduct struct {
	ctrl     *gomock.Controller
	recorder *MockProductMockRecorder
	isgomock struct{}
}

// MockProductMockRecorder is the mock recorder for MockProduct.
type MockProductMockRecorder struct {
	mock *MockProduct
}

// NewMockProduct creates a new mock instance.
func NewMockProduct(ctrl *gomock.Controller) *MockProduct {
	mock := &MockProduct{ctrl: ctrl}
	mock.recorder = &MockProductMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProduct) EXPECT() *MockProductMockRecorder {
	return m.recorder
}

// AddProduct mocks base method.
func (m *MockProduct) AddProduct(ctx context.Context, pvzId uuid.UUID, productType string) (model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddProduct", ctx, pvzId, productType)
	ret0, _ := ret[0].(model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddProduct indicates an expected call of AddProduct.
func (mr *MockProductMockRecorder) AddProduct(ctx, pvzId, productType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddProduct", reflect.TypeOf((*MockProduct)(nil).AddProduct), ctx, pvzId, productType)
}

// DeleteLastProduct mocks base method.
func (m *MockProduct) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastProduct", ctx, pvzId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastProduct indicates an expected call of DeleteLastProduct.
func (mr *MockProductMockRecorder) DeleteLastProduct(ctx, pvzId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastProduct", reflect.TypeOf((*MockProduct)(nil).DeleteLastProduct), ctx, pvzId)
}
