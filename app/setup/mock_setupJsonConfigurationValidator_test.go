// Code generated by mockery v2.20.0. DO NOT EDIT.

package setup

import (
	context "context"

	appcontext "github.com/cloudogu/k8s-ces-setup/app/context"

	mock "github.com/stretchr/testify/mock"
)

// mockSetupJsonConfigurationValidator is an autogenerated mock type for the setupJsonConfigurationValidator type
type mockSetupJsonConfigurationValidator struct {
	mock.Mock
}

type mockSetupJsonConfigurationValidator_Expecter struct {
	mock *mock.Mock
}

func (_m *mockSetupJsonConfigurationValidator) EXPECT() *mockSetupJsonConfigurationValidator_Expecter {
	return &mockSetupJsonConfigurationValidator_Expecter{mock: &_m.Mock}
}

// Validate provides a mock function with given fields: ctx, setupJson
func (_m *mockSetupJsonConfigurationValidator) Validate(ctx context.Context, setupJson *appcontext.SetupJsonConfiguration) error {
	ret := _m.Called(ctx, setupJson)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *appcontext.SetupJsonConfiguration) error); ok {
		r0 = rf(ctx, setupJson)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSetupJsonConfigurationValidator_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type mockSetupJsonConfigurationValidator_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - ctx context.Context
//   - setupJson *appcontext.SetupJsonConfiguration
func (_e *mockSetupJsonConfigurationValidator_Expecter) Validate(ctx interface{}, setupJson interface{}) *mockSetupJsonConfigurationValidator_Validate_Call {
	return &mockSetupJsonConfigurationValidator_Validate_Call{Call: _e.mock.On("Validate", ctx, setupJson)}
}

func (_c *mockSetupJsonConfigurationValidator_Validate_Call) Run(run func(ctx context.Context, setupJson *appcontext.SetupJsonConfiguration)) *mockSetupJsonConfigurationValidator_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*appcontext.SetupJsonConfiguration))
	})
	return _c
}

func (_c *mockSetupJsonConfigurationValidator_Validate_Call) Return(_a0 error) *mockSetupJsonConfigurationValidator_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSetupJsonConfigurationValidator_Validate_Call) RunAndReturn(run func(context.Context, *appcontext.SetupJsonConfiguration) error) *mockSetupJsonConfigurationValidator_Validate_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockSetupJsonConfigurationValidator interface {
	mock.TestingT
	Cleanup(func())
}

// newMockSetupJsonConfigurationValidator creates a new instance of mockSetupJsonConfigurationValidator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockSetupJsonConfigurationValidator(t mockConstructorTestingTnewMockSetupJsonConfigurationValidator) *mockSetupJsonConfigurationValidator {
	mock := &mockSetupJsonConfigurationValidator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
