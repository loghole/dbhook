// Code generated by MockGen. DO NOT EDIT.
// Source: tests/interfaces/conn.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	driver "database/sql/driver"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAdvancedConn is a mock of AdvancedConn interface.
type MockAdvancedConn struct {
	ctrl     *gomock.Controller
	recorder *MockAdvancedConnMockRecorder
}

// MockAdvancedConnMockRecorder is the mock recorder for MockAdvancedConn.
type MockAdvancedConnMockRecorder struct {
	mock *MockAdvancedConn
}

// NewMockAdvancedConn creates a new mock instance.
func NewMockAdvancedConn(ctrl *gomock.Controller) *MockAdvancedConn {
	mock := &MockAdvancedConn{ctrl: ctrl}
	mock.recorder = &MockAdvancedConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdvancedConn) EXPECT() *MockAdvancedConnMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockAdvancedConn) Begin() (driver.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(driver.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockAdvancedConnMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockAdvancedConn)(nil).Begin))
}

// BeginTx mocks base method.
func (m *MockAdvancedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, opts)
	ret0, _ := ret[0].(driver.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockAdvancedConnMockRecorder) BeginTx(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockAdvancedConn)(nil).BeginTx), ctx, opts)
}

// CheckNamedValue mocks base method.
func (m *MockAdvancedConn) CheckNamedValue(arg0 *driver.NamedValue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckNamedValue", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckNamedValue indicates an expected call of CheckNamedValue.
func (mr *MockAdvancedConnMockRecorder) CheckNamedValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckNamedValue", reflect.TypeOf((*MockAdvancedConn)(nil).CheckNamedValue), arg0)
}

// Close mocks base method.
func (m *MockAdvancedConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockAdvancedConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAdvancedConn)(nil).Close))
}

// ExecContext mocks base method.
func (m *MockAdvancedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecContext", ctx, query, args)
	ret0, _ := ret[0].(driver.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockAdvancedConnMockRecorder) ExecContext(ctx, query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockAdvancedConn)(nil).ExecContext), ctx, query, args)
}

// Prepare mocks base method.
func (m *MockAdvancedConn) Prepare(query string) (driver.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(driver.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockAdvancedConnMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockAdvancedConn)(nil).Prepare), query)
}

// PrepareContext mocks base method.
func (m *MockAdvancedConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareContext", ctx, query)
	ret0, _ := ret[0].(driver.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareContext indicates an expected call of PrepareContext.
func (mr *MockAdvancedConnMockRecorder) PrepareContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareContext", reflect.TypeOf((*MockAdvancedConn)(nil).PrepareContext), ctx, query)
}

// QueryContext mocks base method.
func (m *MockAdvancedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryContext", ctx, query, args)
	ret0, _ := ret[0].(driver.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockAdvancedConnMockRecorder) QueryContext(ctx, query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockAdvancedConn)(nil).QueryContext), ctx, query, args)
}

// ResetSession mocks base method.
func (m *MockAdvancedConn) ResetSession(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetSession", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetSession indicates an expected call of ResetSession.
func (mr *MockAdvancedConnMockRecorder) ResetSession(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetSession", reflect.TypeOf((*MockAdvancedConn)(nil).ResetSession), ctx)
}

// MockQueryerConn is a mock of QueryerConn interface.
type MockQueryerConn struct {
	ctrl     *gomock.Controller
	recorder *MockQueryerConnMockRecorder
}

// MockQueryerConnMockRecorder is the mock recorder for MockQueryerConn.
type MockQueryerConnMockRecorder struct {
	mock *MockQueryerConn
}

// NewMockQueryerConn creates a new mock instance.
func NewMockQueryerConn(ctrl *gomock.Controller) *MockQueryerConn {
	mock := &MockQueryerConn{ctrl: ctrl}
	mock.recorder = &MockQueryerConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryerConn) EXPECT() *MockQueryerConnMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockQueryerConn) Begin() (driver.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(driver.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockQueryerConnMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockQueryerConn)(nil).Begin))
}

// Close mocks base method.
func (m *MockQueryerConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockQueryerConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockQueryerConn)(nil).Close))
}

// Exec mocks base method.
func (m *MockQueryerConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", query, args)
	ret0, _ := ret[0].(driver.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockQueryerConnMockRecorder) Exec(query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockQueryerConn)(nil).Exec), query, args)
}

// Prepare mocks base method.
func (m *MockQueryerConn) Prepare(query string) (driver.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(driver.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockQueryerConnMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockQueryerConn)(nil).Prepare), query)
}

// Query mocks base method.
func (m *MockQueryerConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", query, args)
	ret0, _ := ret[0].(driver.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockQueryerConnMockRecorder) Query(query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockQueryerConn)(nil).Query), query, args)
}

// ResetSession mocks base method.
func (m *MockQueryerConn) ResetSession(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetSession", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetSession indicates an expected call of ResetSession.
func (mr *MockQueryerConnMockRecorder) ResetSession(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetSession", reflect.TypeOf((*MockQueryerConn)(nil).ResetSession), ctx)
}

// MockQConn is a mock of QConn interface.
type MockQConn struct {
	ctrl     *gomock.Controller
	recorder *MockQConnMockRecorder
}

// MockQConnMockRecorder is the mock recorder for MockQConn.
type MockQConnMockRecorder struct {
	mock *MockQConn
}

// NewMockQConn creates a new mock instance.
func NewMockQConn(ctrl *gomock.Controller) *MockQConn {
	mock := &MockQConn{ctrl: ctrl}
	mock.recorder = &MockQConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQConn) EXPECT() *MockQConnMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockQConn) Begin() (driver.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(driver.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockQConnMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockQConn)(nil).Begin))
}

// Close mocks base method.
func (m *MockQConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockQConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockQConn)(nil).Close))
}

// Prepare mocks base method.
func (m *MockQConn) Prepare(query string) (driver.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(driver.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockQConnMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockQConn)(nil).Prepare), query)
}

// QueryContext mocks base method.
func (m *MockQConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryContext", ctx, query, args)
	ret0, _ := ret[0].(driver.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockQConnMockRecorder) QueryContext(ctx, query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockQConn)(nil).QueryContext), ctx, query, args)
}

// MockEConn is a mock of EConn interface.
type MockEConn struct {
	ctrl     *gomock.Controller
	recorder *MockEConnMockRecorder
}

// MockEConnMockRecorder is the mock recorder for MockEConn.
type MockEConnMockRecorder struct {
	mock *MockEConn
}

// NewMockEConn creates a new mock instance.
func NewMockEConn(ctrl *gomock.Controller) *MockEConn {
	mock := &MockEConn{ctrl: ctrl}
	mock.recorder = &MockEConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEConn) EXPECT() *MockEConnMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockEConn) Begin() (driver.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin")
	ret0, _ := ret[0].(driver.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockEConnMockRecorder) Begin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockEConn)(nil).Begin))
}

// Close mocks base method.
func (m *MockEConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockEConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockEConn)(nil).Close))
}

// ExecContext mocks base method.
func (m *MockEConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecContext", ctx, query, args)
	ret0, _ := ret[0].(driver.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockEConnMockRecorder) ExecContext(ctx, query, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockEConn)(nil).ExecContext), ctx, query, args)
}

// Prepare mocks base method.
func (m *MockEConn) Prepare(query string) (driver.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", query)
	ret0, _ := ret[0].(driver.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockEConnMockRecorder) Prepare(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockEConn)(nil).Prepare), query)
}
