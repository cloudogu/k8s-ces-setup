// Code generated by mockery v2.42.1. DO NOT EDIT.

package data

import mock "github.com/stretchr/testify/mock"

// MockInternalConfigRegistryProvider is an autogenerated mock type for the InternalConfigRegistryProvider type
type MockInternalConfigRegistryProvider struct {
	mock.Mock
}

type MockInternalConfigRegistryProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInternalConfigRegistryProvider) EXPECT() *MockInternalConfigRegistryProvider_Expecter {
	return &MockInternalConfigRegistryProvider_Expecter{mock: &_m.Mock}
}

// NewMockInternalConfigRegistryProvider creates a new instance of MockInternalConfigRegistryProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInternalConfigRegistryProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInternalConfigRegistryProvider {
	mock := &MockInternalConfigRegistryProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
