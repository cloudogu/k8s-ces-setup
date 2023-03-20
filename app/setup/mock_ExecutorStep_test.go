// Code generated by mockery v2.20.0. DO NOT EDIT.

package setup

import mock "github.com/stretchr/testify/mock"

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

// GetStepDescription provides a mock function with given fields:
func (_m *MockExecutorStep) GetStepDescription() string {
	ret := _m.Called()

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

// PerformSetupStep provides a mock function with given fields:
func (_m *MockExecutorStep) PerformSetupStep() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
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
func (_e *MockExecutorStep_Expecter) PerformSetupStep() *MockExecutorStep_PerformSetupStep_Call {
	return &MockExecutorStep_PerformSetupStep_Call{Call: _e.mock.On("PerformSetupStep")}
}

func (_c *MockExecutorStep_PerformSetupStep_Call) Run(run func()) *MockExecutorStep_PerformSetupStep_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockExecutorStep_PerformSetupStep_Call) Return(_a0 error) *MockExecutorStep_PerformSetupStep_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockExecutorStep_PerformSetupStep_Call) RunAndReturn(run func() error) *MockExecutorStep_PerformSetupStep_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockExecutorStep interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockExecutorStep creates a new instance of MockExecutorStep. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockExecutorStep(t mockConstructorTestingTNewMockExecutorStep) *MockExecutorStep {
	mock := &MockExecutorStep{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}