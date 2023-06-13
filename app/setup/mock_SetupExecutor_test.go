// Code generated by mockery v2.20.0. DO NOT EDIT.

package setup

import (
	registry "github.com/cloudogu/cesapp-lib/registry"
	mock "github.com/stretchr/testify/mock"
)

// MockSetupExecutor is an autogenerated mock type for the SetupExecutor type
type MockSetupExecutor struct {
	mock.Mock
}

type MockSetupExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSetupExecutor) EXPECT() *MockSetupExecutor_Expecter {
	return &MockSetupExecutor_Expecter{mock: &_m.Mock}
}

// PerformSetup provides a mock function with given fields:
func (_m *MockSetupExecutor) PerformSetup() (error, string) {
	ret := _m.Called()

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func() (error, string)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// MockSetupExecutor_PerformSetup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PerformSetup'
type MockSetupExecutor_PerformSetup_Call struct {
	*mock.Call
}

// PerformSetup is a helper method to define mock.On call
func (_e *MockSetupExecutor_Expecter) PerformSetup() *MockSetupExecutor_PerformSetup_Call {
	return &MockSetupExecutor_PerformSetup_Call{Call: _e.mock.On("PerformSetup")}
}

func (_c *MockSetupExecutor_PerformSetup_Call) Run(run func()) *MockSetupExecutor_PerformSetup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSetupExecutor_PerformSetup_Call) Return(_a0 error, _a1 string) *MockSetupExecutor_PerformSetup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSetupExecutor_PerformSetup_Call) RunAndReturn(run func() (error, string)) *MockSetupExecutor_PerformSetup_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterComponentSetupSteps provides a mock function with given fields:
func (_m *MockSetupExecutor) RegisterComponentSetupSteps() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSetupExecutor_RegisterComponentSetupSteps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterComponentSetupSteps'
type MockSetupExecutor_RegisterComponentSetupSteps_Call struct {
	*mock.Call
}

// RegisterComponentSetupSteps is a helper method to define mock.On call
func (_e *MockSetupExecutor_Expecter) RegisterComponentSetupSteps() *MockSetupExecutor_RegisterComponentSetupSteps_Call {
	return &MockSetupExecutor_RegisterComponentSetupSteps_Call{Call: _e.mock.On("RegisterComponentSetupSteps")}
}

func (_c *MockSetupExecutor_RegisterComponentSetupSteps_Call) Run(run func()) *MockSetupExecutor_RegisterComponentSetupSteps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSetupExecutor_RegisterComponentSetupSteps_Call) Return(_a0 error) *MockSetupExecutor_RegisterComponentSetupSteps_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSetupExecutor_RegisterComponentSetupSteps_Call) RunAndReturn(run func() error) *MockSetupExecutor_RegisterComponentSetupSteps_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterDataSetupSteps provides a mock function with given fields: _a0
func (_m *MockSetupExecutor) RegisterDataSetupSteps(_a0 registry.Registry) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(registry.Registry) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSetupExecutor_RegisterDataSetupSteps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterDataSetupSteps'
type MockSetupExecutor_RegisterDataSetupSteps_Call struct {
	*mock.Call
}

// RegisterDataSetupSteps is a helper method to define mock.On call
//   - _a0 registry.Registry
func (_e *MockSetupExecutor_Expecter) RegisterDataSetupSteps(_a0 interface{}) *MockSetupExecutor_RegisterDataSetupSteps_Call {
	return &MockSetupExecutor_RegisterDataSetupSteps_Call{Call: _e.mock.On("RegisterDataSetupSteps", _a0)}
}

func (_c *MockSetupExecutor_RegisterDataSetupSteps_Call) Run(run func(_a0 registry.Registry)) *MockSetupExecutor_RegisterDataSetupSteps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(registry.Registry))
	})
	return _c
}

func (_c *MockSetupExecutor_RegisterDataSetupSteps_Call) Return(_a0 error) *MockSetupExecutor_RegisterDataSetupSteps_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSetupExecutor_RegisterDataSetupSteps_Call) RunAndReturn(run func(registry.Registry) error) *MockSetupExecutor_RegisterDataSetupSteps_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterDoguInstallationSteps provides a mock function with given fields:
func (_m *MockSetupExecutor) RegisterDoguInstallationSteps() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSetupExecutor_RegisterDoguInstallationSteps_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterDoguInstallationSteps'
type MockSetupExecutor_RegisterDoguInstallationSteps_Call struct {
	*mock.Call
}

// RegisterDoguInstallationSteps is a helper method to define mock.On call
func (_e *MockSetupExecutor_Expecter) RegisterDoguInstallationSteps() *MockSetupExecutor_RegisterDoguInstallationSteps_Call {
	return &MockSetupExecutor_RegisterDoguInstallationSteps_Call{Call: _e.mock.On("RegisterDoguInstallationSteps")}
}

func (_c *MockSetupExecutor_RegisterDoguInstallationSteps_Call) Run(run func()) *MockSetupExecutor_RegisterDoguInstallationSteps_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSetupExecutor_RegisterDoguInstallationSteps_Call) Return(_a0 error) *MockSetupExecutor_RegisterDoguInstallationSteps_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSetupExecutor_RegisterDoguInstallationSteps_Call) RunAndReturn(run func() error) *MockSetupExecutor_RegisterDoguInstallationSteps_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterFQDNRetrieverStep provides a mock function with given fields:
func (_m *MockSetupExecutor) RegisterFQDNRetrieverStep() {
	_m.Called()
}

// MockSetupExecutor_RegisterFQDNRetrieverStep_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterFQDNRetrieverStep'
type MockSetupExecutor_RegisterFQDNRetrieverStep_Call struct {
	*mock.Call
}

// RegisterFQDNRetrieverStep is a helper method to define mock.On call
func (_e *MockSetupExecutor_Expecter) RegisterFQDNRetrieverStep() *MockSetupExecutor_RegisterFQDNRetrieverStep_Call {
	return &MockSetupExecutor_RegisterFQDNRetrieverStep_Call{Call: _e.mock.On("RegisterFQDNRetrieverStep")}
}

func (_c *MockSetupExecutor_RegisterFQDNRetrieverStep_Call) Run(run func()) *MockSetupExecutor_RegisterFQDNRetrieverStep_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSetupExecutor_RegisterFQDNRetrieverStep_Call) Return() *MockSetupExecutor_RegisterFQDNRetrieverStep_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockSetupExecutor_RegisterFQDNRetrieverStep_Call) RunAndReturn(run func()) *MockSetupExecutor_RegisterFQDNRetrieverStep_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterSSLGenerationStep provides a mock function with given fields:
func (_m *MockSetupExecutor) RegisterSSLGenerationStep() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSetupExecutor_RegisterSSLGenerationStep_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterSSLGenerationStep'
type MockSetupExecutor_RegisterSSLGenerationStep_Call struct {
	*mock.Call
}

// RegisterSSLGenerationStep is a helper method to define mock.On call
func (_e *MockSetupExecutor_Expecter) RegisterSSLGenerationStep() *MockSetupExecutor_RegisterSSLGenerationStep_Call {
	return &MockSetupExecutor_RegisterSSLGenerationStep_Call{Call: _e.mock.On("RegisterSSLGenerationStep")}
}

func (_c *MockSetupExecutor_RegisterSSLGenerationStep_Call) Run(run func()) *MockSetupExecutor_RegisterSSLGenerationStep_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSetupExecutor_RegisterSSLGenerationStep_Call) Return(_a0 error) *MockSetupExecutor_RegisterSSLGenerationStep_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSetupExecutor_RegisterSSLGenerationStep_Call) RunAndReturn(run func() error) *MockSetupExecutor_RegisterSSLGenerationStep_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterValidationStep provides a mock function with given fields:
func (_m *MockSetupExecutor) RegisterValidationStep() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSetupExecutor_RegisterValidationStep_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterValidationStep'
type MockSetupExecutor_RegisterValidationStep_Call struct {
	*mock.Call
}

// RegisterValidationStep is a helper method to define mock.On call
func (_e *MockSetupExecutor_Expecter) RegisterValidationStep() *MockSetupExecutor_RegisterValidationStep_Call {
	return &MockSetupExecutor_RegisterValidationStep_Call{Call: _e.mock.On("RegisterValidationStep")}
}

func (_c *MockSetupExecutor_RegisterValidationStep_Call) Run(run func()) *MockSetupExecutor_RegisterValidationStep_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSetupExecutor_RegisterValidationStep_Call) Return(_a0 error) *MockSetupExecutor_RegisterValidationStep_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSetupExecutor_RegisterValidationStep_Call) RunAndReturn(run func() error) *MockSetupExecutor_RegisterValidationStep_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockSetupExecutor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSetupExecutor creates a new instance of MockSetupExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSetupExecutor(t mockConstructorTestingTNewMockSetupExecutor) *MockSetupExecutor {
	mock := &MockSetupExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
