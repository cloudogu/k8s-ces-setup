// Code generated by mockery v2.53.3. DO NOT EDIT.

package setup

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockExecutorStep is an autogenerated mock type for the ExecutorStep type
type MockExecutorStep struct {
	mock.Mock
}

type MockExecutorStep_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExecutorStep) EXPECT() *MockExecutorStep_Expecter {
	return &MockExecutorStep_Expecter{mock: &_m.Mock}
}

// GetStepDescription provides a mock function with no fields
func (_m *MockExecutorStep) GetStepDescription() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetStepDescription")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockExecutorStep_GetStepDescription_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStepDescription'
type MockExecutorStep_GetStepDescription_Call struct {
	*mock.Call
}

// GetStepDescription is a helper method to define mock.On call
func (_e *MockExecutorStep_Expecter) GetStepDescription() *MockExecutorStep_GetStepDescription_Call {
	return &MockExecutorStep_GetStepDescription_Call{Call: _e.mock.On("GetStepDescription")}
}

func (_c *MockExecutorStep_GetStepDescription_Call) Run(run func()) *MockExecutorStep_GetStepDescription_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExecutorStep_GetStepDescription_Call) Return(_a0 string) *MockExecutorStep_GetStepDescription_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExecutorStep_GetStepDescription_Call) RunAndReturn(run func() string) *MockExecutorStep_GetStepDescription_Call {
	_c.Call.Return(run)
	return _c
}

// PerformSetupStep provides a mock function with given fields: ctx
func (_m *MockExecutorStep) PerformSetupStep(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for PerformSetupStep")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExecutorStep_PerformSetupStep_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PerformSetupStep'
type MockExecutorStep_PerformSetupStep_Call struct {
	*mock.Call
}

// PerformSetupStep is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockExecutorStep_Expecter) PerformSetupStep(ctx interface{}) *MockExecutorStep_PerformSetupStep_Call {
	return &MockExecutorStep_PerformSetupStep_Call{Call: _e.mock.On("PerformSetupStep", ctx)}
}

func (_c *MockExecutorStep_PerformSetupStep_Call) Run(run func(ctx context.Context)) *MockExecutorStep_PerformSetupStep_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockExecutorStep_PerformSetupStep_Call) Return(_a0 error) *MockExecutorStep_PerformSetupStep_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExecutorStep_PerformSetupStep_Call) RunAndReturn(run func(context.Context) error) *MockExecutorStep_PerformSetupStep_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExecutorStep creates a new instance of MockExecutorStep. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExecutorStep(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExecutorStep {
	mock := &MockExecutorStep{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
