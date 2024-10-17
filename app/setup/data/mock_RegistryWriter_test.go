// Code generated by mockery v2.20.0. DO NOT EDIT.

package data

import (
	context "github.com/cloudogu/k8s-ces-setup/app/context"
	mock "github.com/stretchr/testify/mock"
)

// MockRegistryWriter is an autogenerated mock type for the RegistryWriter type
type MockRegistryWriter struct {
	mock.Mock
}

type MockRegistryWriter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRegistryWriter) EXPECT() *MockRegistryWriter_Expecter {
	return &MockRegistryWriter_Expecter{mock: &_m.Mock}
}

// WriteConfigToRegistry provides a mock function with given fields: registryConfig
func (_m *MockRegistryWriter) WriteConfigToRegistry(registryConfig context.CustomKeyValue) error {
	ret := _m.Called(registryConfig)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.CustomKeyValue) error); ok {
		r0 = rf(registryConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRegistryWriter_WriteConfigToRegistry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteConfigToRegistry'
type MockRegistryWriter_WriteConfigToRegistry_Call struct {
	*mock.Call
}

// WriteConfigToRegistry is a helper method to define mock.On call
//   - registryConfig context.CustomKeyValue
func (_e *MockRegistryWriter_Expecter) WriteConfigToRegistry(registryConfig interface{}) *MockRegistryWriter_WriteConfigToRegistry_Call {
	return &MockRegistryWriter_WriteConfigToRegistry_Call{Call: _e.mock.On("WriteConfigToRegistry", registryConfig)}
}

func (_c *MockRegistryWriter_WriteConfigToRegistry_Call) Run(run func(registryConfig context.CustomKeyValue)) *MockRegistryWriter_WriteConfigToRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.CustomKeyValue))
	})
	return _c
}

func (_c *MockRegistryWriter_WriteConfigToRegistry_Call) Return(_a0 error) *MockRegistryWriter_WriteConfigToRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRegistryWriter_WriteConfigToRegistry_Call) RunAndReturn(run func(context.CustomKeyValue) error) *MockRegistryWriter_WriteConfigToRegistry_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockRegistryWriter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRegistryWriter creates a new instance of MockRegistryWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRegistryWriter(t mockConstructorTestingTNewMockRegistryWriter) *MockRegistryWriter {
	mock := &MockRegistryWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
