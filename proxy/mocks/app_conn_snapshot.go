// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/airchains-network/tracksbft/abci/types"
)

// AppConnSnapshot is an autogenerated mock type for the AppConnSnapshot type
type AppConnSnapshot struct {
	mock.Mock
}

// ApplySnapshotChunkSync provides a mock function with given fields: _a0
func (_m *AppConnSnapshot) ApplySnapshotChunkSync(_a0 types.RequestApplySnapshotChunk) (*types.ResponseApplySnapshotChunk, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseApplySnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestApplySnapshotChunk) *types.ResponseApplySnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseApplySnapshotChunk)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestApplySnapshotChunk) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Error provides a mock function with given fields:
func (_m *AppConnSnapshot) Error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListSnapshotsSync provides a mock function with given fields: _a0
func (_m *AppConnSnapshot) ListSnapshotsSync(_a0 types.RequestListSnapshots) (*types.ResponseListSnapshots, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseListSnapshots
	if rf, ok := ret.Get(0).(func(types.RequestListSnapshots) *types.ResponseListSnapshots); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseListSnapshots)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestListSnapshots) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadSnapshotChunkSync provides a mock function with given fields: _a0
func (_m *AppConnSnapshot) LoadSnapshotChunkSync(_a0 types.RequestLoadSnapshotChunk) (*types.ResponseLoadSnapshotChunk, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseLoadSnapshotChunk
	if rf, ok := ret.Get(0).(func(types.RequestLoadSnapshotChunk) *types.ResponseLoadSnapshotChunk); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseLoadSnapshotChunk)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestLoadSnapshotChunk) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OfferSnapshotSync provides a mock function with given fields: _a0
func (_m *AppConnSnapshot) OfferSnapshotSync(_a0 types.RequestOfferSnapshot) (*types.ResponseOfferSnapshot, error) {
	ret := _m.Called(_a0)

	var r0 *types.ResponseOfferSnapshot
	if rf, ok := ret.Get(0).(func(types.RequestOfferSnapshot) *types.ResponseOfferSnapshot); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.ResponseOfferSnapshot)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.RequestOfferSnapshot) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAppConnSnapshot interface {
	mock.TestingT
	Cleanup(func())
}

// NewAppConnSnapshot creates a new instance of AppConnSnapshot. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAppConnSnapshot(t mockConstructorTestingTNewAppConnSnapshot) *AppConnSnapshot {
	mock := &AppConnSnapshot{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
