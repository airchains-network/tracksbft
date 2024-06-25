// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	abcicli "github.com/airchains-network/tracksbft/abci/client"
	mock "github.com/stretchr/testify/mock"
)

// ClientCreator is an autogenerated mock type for the ClientCreator type
type ClientCreator struct {
	mock.Mock
}

// NewABCIClient provides a mock function with given fields:
func (_m *ClientCreator) NewABCIClient() (abcicli.Client, error) {
	ret := _m.Called()

	var r0 abcicli.Client
	if rf, ok := ret.Get(0).(func() abcicli.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(abcicli.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewClientCreator interface {
	mock.TestingT
	Cleanup(func())
}

// NewClientCreator creates a new instance of ClientCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClientCreator(t mockConstructorTestingTNewClientCreator) *ClientCreator {
	mock := &ClientCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}