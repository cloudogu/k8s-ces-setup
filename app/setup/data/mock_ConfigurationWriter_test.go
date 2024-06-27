// Code generated by mockery v2.42.1. DO NOT EDIT.

package data

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockConfigurationWriter is an autogenerated mock type for the ConfigurationRegistry type
type MockConfigurationWriter struct {
	mock.Mock
}

type MockConfigurationWriter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConfigurationWriter) EXPECT() *MockConfigurationWriter_Expecter {
	return &MockConfigurationWriter_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, key
func (_m *MockConfigurationWriter) Delete(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConfigurationWriter_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockConfigurationWriter_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockConfigurationWriter_Expecter) Delete(ctx interface{}, key interface{}) *MockConfigurationWriter_Delete_Call {
	return &MockConfigurationWriter_Delete_Call{Call: _e.mock.On("Delete", ctx, key)}
}

func (_c *MockConfigurationWriter_Delete_Call) Run(run func(ctx context.Context, key string)) *MockConfigurationWriter_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockConfigurationWriter_Delete_Call) Return(_a0 error) *MockConfigurationWriter_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfigurationWriter_Delete_Call) RunAndReturn(run func(context.Context, string) error) *MockConfigurationWriter_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAll provides a mock function with given fields: ctx
func (_m *MockConfigurationWriter) DeleteAll(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConfigurationWriter_DeleteAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAll'
type MockConfigurationWriter_DeleteAll_Call struct {
	*mock.Call
}

// DeleteAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockConfigurationWriter_Expecter) DeleteAll(ctx interface{}) *MockConfigurationWriter_DeleteAll_Call {
	return &MockConfigurationWriter_DeleteAll_Call{Call: _e.mock.On("DeleteAll", ctx)}
}

func (_c *MockConfigurationWriter_DeleteAll_Call) Run(run func(ctx context.Context)) *MockConfigurationWriter_DeleteAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockConfigurationWriter_DeleteAll_Call) Return(_a0 error) *MockConfigurationWriter_DeleteAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfigurationWriter_DeleteAll_Call) RunAndReturn(run func(context.Context) error) *MockConfigurationWriter_DeleteAll_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteRecursive provides a mock function with given fields: ctx, key
func (_m *MockConfigurationWriter) DeleteRecursive(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for DeleteRecursive")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConfigurationWriter_DeleteRecursive_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteRecursive'
type MockConfigurationWriter_DeleteRecursive_Call struct {
	*mock.Call
}

// DeleteRecursive is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *MockConfigurationWriter_Expecter) DeleteRecursive(ctx interface{}, key interface{}) *MockConfigurationWriter_DeleteRecursive_Call {
	return &MockConfigurationWriter_DeleteRecursive_Call{Call: _e.mock.On("DeleteRecursive", ctx, key)}
}

func (_c *MockConfigurationWriter_DeleteRecursive_Call) Run(run func(ctx context.Context, key string)) *MockConfigurationWriter_DeleteRecursive_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockConfigurationWriter_DeleteRecursive_Call) Return(_a0 error) *MockConfigurationWriter_DeleteRecursive_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfigurationWriter_DeleteRecursive_Call) RunAndReturn(run func(context.Context, string) error) *MockConfigurationWriter_DeleteRecursive_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value
func (_m *MockConfigurationWriter) Set(ctx context.Context, key string, value string) error {
	ret := _m.Called(ctx, key, value)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockConfigurationWriter_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type MockConfigurationWriter_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value string
func (_e *MockConfigurationWriter_Expecter) Set(ctx interface{}, key interface{}, value interface{}) *MockConfigurationWriter_Set_Call {
	return &MockConfigurationWriter_Set_Call{Call: _e.mock.On("Set", ctx, key, value)}
}

func (_c *MockConfigurationWriter_Set_Call) Run(run func(ctx context.Context, key string, value string)) *MockConfigurationWriter_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockConfigurationWriter_Set_Call) Return(_a0 error) *MockConfigurationWriter_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockConfigurationWriter_Set_Call) RunAndReturn(run func(context.Context, string, string) error) *MockConfigurationWriter_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConfigurationWriter creates a new instance of MockConfigurationWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConfigurationWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConfigurationWriter {
	mock := &MockConfigurationWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
