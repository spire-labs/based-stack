// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	eth "github.com/ethereum-optimism/optimism/op-service/eth"
	mock "github.com/stretchr/testify/mock"
)

// BeaconClient is an autogenerated mock type for the BeaconClient type
type BeaconClient struct {
	mock.Mock
}

type BeaconClient_Expecter struct {
	mock *mock.Mock
}

func (_m *BeaconClient) EXPECT() *BeaconClient_Expecter {
	return &BeaconClient_Expecter{mock: &_m.Mock}
}

// BeaconBlobSideCars provides a mock function with given fields: ctx, fetchAllSidecars, slot, hashes
func (_m *BeaconClient) BeaconBlobSideCars(ctx context.Context, fetchAllSidecars bool, slot uint64, hashes []eth.IndexedBlobHash) (eth.APIGetBlobSidecarsResponse, error) {
	ret := _m.Called(ctx, fetchAllSidecars, slot, hashes)

	var r0 eth.APIGetBlobSidecarsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, uint64, []eth.IndexedBlobHash) (eth.APIGetBlobSidecarsResponse, error)); ok {
		return rf(ctx, fetchAllSidecars, slot, hashes)
	}
	if rf, ok := ret.Get(0).(func(context.Context, bool, uint64, []eth.IndexedBlobHash) eth.APIGetBlobSidecarsResponse); ok {
		r0 = rf(ctx, fetchAllSidecars, slot, hashes)
	} else {
		r0 = ret.Get(0).(eth.APIGetBlobSidecarsResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, bool, uint64, []eth.IndexedBlobHash) error); ok {
		r1 = rf(ctx, fetchAllSidecars, slot, hashes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconClient_BeaconBlobSideCars_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BeaconBlobSideCars'
type BeaconClient_BeaconBlobSideCars_Call struct {
	*mock.Call
}

// BeaconBlobSideCars is a helper method to define mock.On call
//   - ctx context.Context
//   - fetchAllSidecars bool
//   - slot uint64
//   - hashes []eth.IndexedBlobHash
func (_e *BeaconClient_Expecter) BeaconBlobSideCars(ctx interface{}, fetchAllSidecars interface{}, slot interface{}, hashes interface{}) *BeaconClient_BeaconBlobSideCars_Call {
	return &BeaconClient_BeaconBlobSideCars_Call{Call: _e.mock.On("BeaconBlobSideCars", ctx, fetchAllSidecars, slot, hashes)}
}

func (_c *BeaconClient_BeaconBlobSideCars_Call) Run(run func(ctx context.Context, fetchAllSidecars bool, slot uint64, hashes []eth.IndexedBlobHash)) *BeaconClient_BeaconBlobSideCars_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(bool), args[2].(uint64), args[3].([]eth.IndexedBlobHash))
	})
	return _c
}

func (_c *BeaconClient_BeaconBlobSideCars_Call) Return(_a0 eth.APIGetBlobSidecarsResponse, _a1 error) *BeaconClient_BeaconBlobSideCars_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BeaconClient_BeaconBlobSideCars_Call) RunAndReturn(run func(context.Context, bool, uint64, []eth.IndexedBlobHash) (eth.APIGetBlobSidecarsResponse, error)) *BeaconClient_BeaconBlobSideCars_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *BeaconClient) GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error) {
	ret := _m.Called(ctx, epoch)

	var r0 eth.APIGetLookaheadResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (eth.APIGetLookaheadResponse, error)); ok {
		return rf(ctx, epoch)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) eth.APIGetLookaheadResponse); ok {
		r0 = rf(ctx, epoch)
	} else {
		r0 = ret.Get(0).(eth.APIGetLookaheadResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, epoch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconClient_GetLookahead_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLookahead'
type BeaconClient_GetLookahead_Call struct {
	*mock.Call
}

// BeaconGenesis provides a mock function with given fields: ctx
func (_m *BeaconClient) BeaconGenesis(ctx context.Context) (eth.APIGenesisResponse, error) {
	ret := _m.Called(ctx)

	var r0 eth.APIGenesisResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (eth.APIGenesisResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) eth.APIGenesisResponse); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(eth.APIGenesisResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconClient_BeaconGenesis_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BeaconGenesis'
type BeaconClient_BeaconGenesis_Call struct {
	*mock.Call
}

// BeaconGenesis is a helper method to define mock.On call
//   - ctx context.Context
func (_e *BeaconClient_Expecter) BeaconGenesis(ctx interface{}) *BeaconClient_BeaconGenesis_Call {
	return &BeaconClient_BeaconGenesis_Call{Call: _e.mock.On("BeaconGenesis", ctx)}
}

func (_c *BeaconClient_BeaconGenesis_Call) Run(run func(ctx context.Context)) *BeaconClient_BeaconGenesis_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *BeaconClient_BeaconGenesis_Call) Return(_a0 eth.APIGenesisResponse, _a1 error) *BeaconClient_BeaconGenesis_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BeaconClient_BeaconGenesis_Call) RunAndReturn(run func(context.Context) (eth.APIGenesisResponse, error)) *BeaconClient_BeaconGenesis_Call {
	_c.Call.Return(run)
	return _c
}

// ConfigSpec provides a mock function with given fields: ctx
func (_m *BeaconClient) ConfigSpec(ctx context.Context) (eth.APIConfigResponse, error) {
	ret := _m.Called(ctx)

	var r0 eth.APIConfigResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (eth.APIConfigResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) eth.APIConfigResponse); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(eth.APIConfigResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconClient_ConfigSpec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfigSpec'
type BeaconClient_ConfigSpec_Call struct {
	*mock.Call
}

// ConfigSpec is a helper method to define mock.On call
//   - ctx context.Context
func (_e *BeaconClient_Expecter) ConfigSpec(ctx interface{}) *BeaconClient_ConfigSpec_Call {
	return &BeaconClient_ConfigSpec_Call{Call: _e.mock.On("ConfigSpec", ctx)}
}

func (_c *BeaconClient_ConfigSpec_Call) Run(run func(ctx context.Context)) *BeaconClient_ConfigSpec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *BeaconClient_ConfigSpec_Call) Return(_a0 eth.APIConfigResponse, _a1 error) *BeaconClient_ConfigSpec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BeaconClient_ConfigSpec_Call) RunAndReturn(run func(context.Context) (eth.APIConfigResponse, error)) *BeaconClient_ConfigSpec_Call {
	_c.Call.Return(run)
	return _c
}

// NodeVersion provides a mock function with given fields: ctx
func (_m *BeaconClient) NodeVersion(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeaconClient_NodeVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NodeVersion'
type BeaconClient_NodeVersion_Call struct {
	*mock.Call
}

// NodeVersion is a helper method to define mock.On call
//   - ctx context.Context
func (_e *BeaconClient_Expecter) NodeVersion(ctx interface{}) *BeaconClient_NodeVersion_Call {
	return &BeaconClient_NodeVersion_Call{Call: _e.mock.On("NodeVersion", ctx)}
}

func (_c *BeaconClient_NodeVersion_Call) Run(run func(ctx context.Context)) *BeaconClient_NodeVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *BeaconClient_NodeVersion_Call) Return(_a0 string, _a1 error) *BeaconClient_NodeVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BeaconClient_NodeVersion_Call) RunAndReturn(run func(context.Context) (string, error)) *BeaconClient_NodeVersion_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewBeaconClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewBeaconClient creates a new instance of BeaconClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBeaconClient(t mockConstructorTestingTNewBeaconClient) *BeaconClient {
	mock := &BeaconClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
