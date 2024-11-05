// Code generated by mockery v2.42.1. DO NOT EDIT.

package validation

import (
	context "context"

	dogu "github.com/cloudogu/ces-commons-lib/dogu"
	core "github.com/cloudogu/cesapp-lib/core"

	mock "github.com/stretchr/testify/mock"
)

// mockRemoteDoguDescriptorRepository is an autogenerated mock type for the remoteDoguDescriptorRepository type
type mockRemoteDoguDescriptorRepository struct {
	mock.Mock
}

type mockRemoteDoguDescriptorRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockRemoteDoguDescriptorRepository) EXPECT() *mockRemoteDoguDescriptorRepository_Expecter {
	return &mockRemoteDoguDescriptorRepository_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *mockRemoteDoguDescriptorRepository) Get(_a0 context.Context, _a1 dogu.QualifiedDoguVersion) (*core.Dogu, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.QualifiedDoguVersion) (*core.Dogu, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dogu.QualifiedDoguVersion) *core.Dogu); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dogu.QualifiedDoguVersion) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockRemoteDoguDescriptorRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockRemoteDoguDescriptorRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 dogu.QualifiedDoguVersion
func (_e *mockRemoteDoguDescriptorRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *mockRemoteDoguDescriptorRepository_Get_Call {
	return &mockRemoteDoguDescriptorRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *mockRemoteDoguDescriptorRepository_Get_Call) Run(run func(_a0 context.Context, _a1 dogu.QualifiedDoguVersion)) *mockRemoteDoguDescriptorRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.QualifiedDoguVersion))
	})
	return _c
}

func (_c *mockRemoteDoguDescriptorRepository_Get_Call) Return(_a0 *core.Dogu, _a1 error) *mockRemoteDoguDescriptorRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockRemoteDoguDescriptorRepository_Get_Call) RunAndReturn(run func(context.Context, dogu.QualifiedDoguVersion) (*core.Dogu, error)) *mockRemoteDoguDescriptorRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatest provides a mock function with given fields: _a0, _a1
func (_m *mockRemoteDoguDescriptorRepository) GetLatest(_a0 context.Context, _a1 dogu.QualifiedDoguName) (*core.Dogu, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetLatest")
	}

	var r0 *core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.QualifiedDoguName) (*core.Dogu, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dogu.QualifiedDoguName) *core.Dogu); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dogu.QualifiedDoguName) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockRemoteDoguDescriptorRepository_GetLatest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatest'
type mockRemoteDoguDescriptorRepository_GetLatest_Call struct {
	*mock.Call
}

// GetLatest is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 dogu.QualifiedDoguName
func (_e *mockRemoteDoguDescriptorRepository_Expecter) GetLatest(_a0 interface{}, _a1 interface{}) *mockRemoteDoguDescriptorRepository_GetLatest_Call {
	return &mockRemoteDoguDescriptorRepository_GetLatest_Call{Call: _e.mock.On("GetLatest", _a0, _a1)}
}

func (_c *mockRemoteDoguDescriptorRepository_GetLatest_Call) Run(run func(_a0 context.Context, _a1 dogu.QualifiedDoguName)) *mockRemoteDoguDescriptorRepository_GetLatest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.QualifiedDoguName))
	})
	return _c
}

func (_c *mockRemoteDoguDescriptorRepository_GetLatest_Call) Return(_a0 *core.Dogu, _a1 error) *mockRemoteDoguDescriptorRepository_GetLatest_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockRemoteDoguDescriptorRepository_GetLatest_Call) RunAndReturn(run func(context.Context, dogu.QualifiedDoguName) (*core.Dogu, error)) *mockRemoteDoguDescriptorRepository_GetLatest_Call {
	_c.Call.Return(run)
	return _c
}

// newMockRemoteDoguDescriptorRepository creates a new instance of mockRemoteDoguDescriptorRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRemoteDoguDescriptorRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRemoteDoguDescriptorRepository {
	mock := &mockRemoteDoguDescriptorRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
