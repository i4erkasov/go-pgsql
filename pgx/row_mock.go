// Code generated by mockery v2.40.1. DO NOT EDIT.

package pgx

import mock "github.com/stretchr/testify/mock"

// RowMock is an autogenerated mock type for the RowMock type
type RowMock struct {
	mock.Mock
}

// Scan provides a mock function with given fields: dest
func (_m *RowMock) Scan(dest ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, dest...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Scan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRowMock creates a new instance of RowMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRowMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RowMock {
	mock := &RowMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
