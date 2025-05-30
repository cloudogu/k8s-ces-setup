// Code generated by mockery v2.53.3. DO NOT EDIT.

package setup

import (
	patch "github.com/cloudogu/k8s-ces-setup/v4/app/patch"
	mock "github.com/stretchr/testify/mock"
)

// mockResourcePatchConfigurationValidator is an autogenerated mock type for the resourcePatchConfigurationValidator type
type mockResourcePatchConfigurationValidator struct {
	mock.Mock
}

type mockResourcePatchConfigurationValidator_Expecter struct {
	mock *mock.Mock
}

func (_m *mockResourcePatchConfigurationValidator) EXPECT() *mockResourcePatchConfigurationValidator_Expecter {
	return &mockResourcePatchConfigurationValidator_Expecter{mock: &_m.Mock}
}

// Validate provides a mock function with given fields: resourcePatchConfig
func (_m *mockResourcePatchConfigurationValidator) Validate(resourcePatchConfig []patch.ResourcePatch) error {
	ret := _m.Called(resourcePatchConfig)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]patch.ResourcePatch) error); ok {
		r0 = rf(resourcePatchConfig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockResourcePatchConfigurationValidator_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type mockResourcePatchConfigurationValidator_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - resourcePatchConfig []patch.ResourcePatch
func (_e *mockResourcePatchConfigurationValidator_Expecter) Validate(resourcePatchConfig interface{}) *mockResourcePatchConfigurationValidator_Validate_Call {
	return &mockResourcePatchConfigurationValidator_Validate_Call{Call: _e.mock.On("Validate", resourcePatchConfig)}
}

func (_c *mockResourcePatchConfigurationValidator_Validate_Call) Run(run func(resourcePatchConfig []patch.ResourcePatch)) *mockResourcePatchConfigurationValidator_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]patch.ResourcePatch))
	})
	return _c
}

func (_c *mockResourcePatchConfigurationValidator_Validate_Call) Return(_a0 error) *mockResourcePatchConfigurationValidator_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockResourcePatchConfigurationValidator_Validate_Call) RunAndReturn(run func([]patch.ResourcePatch) error) *mockResourcePatchConfigurationValidator_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// newMockResourcePatchConfigurationValidator creates a new instance of mockResourcePatchConfigurationValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockResourcePatchConfigurationValidator(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockResourcePatchConfigurationValidator {
	mock := &mockResourcePatchConfigurationValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
