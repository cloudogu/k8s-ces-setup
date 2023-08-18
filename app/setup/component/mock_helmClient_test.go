// Code generated by mockery v2.20.0. DO NOT EDIT.

package component

import (
	context "context"

	helmclient "github.com/mittwald/go-helm-client"
	mock "github.com/stretchr/testify/mock"
)

// mockHelmClient is an autogenerated mock type for the helmClient type
type mockHelmClient struct {
	mock.Mock
}

type mockHelmClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockHelmClient) EXPECT() *mockHelmClient_Expecter {
	return &mockHelmClient_Expecter{mock: &_m.Mock}
}

// InstallOrUpgrade provides a mock function with given fields: ctx, chart
func (_m *mockHelmClient) InstallOrUpgrade(ctx context.Context, chart *helmclient.ChartSpec) error {
	ret := _m.Called(ctx, chart)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *helmclient.ChartSpec) error); ok {
		r0 = rf(ctx, chart)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockHelmClient_InstallOrUpgrade_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstallOrUpgrade'
type mockHelmClient_InstallOrUpgrade_Call struct {
	*mock.Call
}

// InstallOrUpgrade is a helper method to define mock.On call
//   - ctx context.Context
//   - chart *helmclient.ChartSpec
func (_e *mockHelmClient_Expecter) InstallOrUpgrade(ctx interface{}, chart interface{}) *mockHelmClient_InstallOrUpgrade_Call {
	return &mockHelmClient_InstallOrUpgrade_Call{Call: _e.mock.On("InstallOrUpgrade", ctx, chart)}
}

func (_c *mockHelmClient_InstallOrUpgrade_Call) Run(run func(ctx context.Context, chart *helmclient.ChartSpec)) *mockHelmClient_InstallOrUpgrade_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*helmclient.ChartSpec))
	})
	return _c
}

func (_c *mockHelmClient_InstallOrUpgrade_Call) Return(_a0 error) *mockHelmClient_InstallOrUpgrade_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockHelmClient_InstallOrUpgrade_Call) RunAndReturn(run func(context.Context, *helmclient.ChartSpec) error) *mockHelmClient_InstallOrUpgrade_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockHelmClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockHelmClient creates a new instance of mockHelmClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockHelmClient(t mockConstructorTestingTnewMockHelmClient) *mockHelmClient {
	mock := &mockHelmClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
