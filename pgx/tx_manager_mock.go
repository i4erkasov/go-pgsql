// Code generated by mockery v2.40.1. DO NOT EDIT.

package pgx

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v4"
)

// TxManagerMock is an autogenerated mock type for the TxManager type
type TxManagerMock struct {
	mock.Mock
}

// WithNestedTx provides a mock function with given fields: ctx, tFunc
func (_m *TxManagerMock) WithNestedTx(ctx context.Context, tFunc func(context.Context, pgx.Tx) error) error {
	ret := _m.Called(ctx, tFunc)

	if len(ret) == 0 {
		panic("no return value specified for WithNestedTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context, pgx.Tx) error) error); ok {
		r0 = rf(ctx, tFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTx provides a mock function with given fields: ctx, fn
func (_m *TxManagerMock) WithTx(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context, pgx.Tx) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTxManagerMock creates a new instance of TxManagerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTxManagerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TxManagerMock {
	mock := &TxManagerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
