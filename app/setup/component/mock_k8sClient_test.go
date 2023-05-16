// Code generated by mockery v2.20.0. DO NOT EDIT.

package component

import (
	apply "github.com/cloudogu/k8s-apply-lib/apply"
	mock "github.com/stretchr/testify/mock"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mockK8sClient is an autogenerated mock type for the k8sClient type
type mockK8sClient struct {
	mock.Mock
}

type mockK8sClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockK8sClient) EXPECT() *mockK8sClient_Expecter {
	return &mockK8sClient_Expecter{mock: &_m.Mock}
}

// Apply provides a mock function with given fields: yamlResources, namespace
func (_m *mockK8sClient) Apply(yamlResources apply.YamlDocument, namespace string) error {
	ret := _m.Called(yamlResources, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(apply.YamlDocument, string) error); ok {
		r0 = rf(yamlResources, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockK8sClient_Apply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Apply'
type mockK8sClient_Apply_Call struct {
	*mock.Call
}

// Apply is a helper method to define mock.On call
//  - yamlResources apply.YamlDocument
//  - namespace string
func (_e *mockK8sClient_Expecter) Apply(yamlResources interface{}, namespace interface{}) *mockK8sClient_Apply_Call {
	return &mockK8sClient_Apply_Call{Call: _e.mock.On("Apply", yamlResources, namespace)}
}

func (_c *mockK8sClient_Apply_Call) Run(run func(yamlResources apply.YamlDocument, namespace string)) *mockK8sClient_Apply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(apply.YamlDocument), args[1].(string))
	})
	return _c
}

func (_c *mockK8sClient_Apply_Call) Return(_a0 error) *mockK8sClient_Apply_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockK8sClient_Apply_Call) RunAndReturn(run func(apply.YamlDocument, string) error) *mockK8sClient_Apply_Call {
	_c.Call.Return(run)
	return _c
}

// ApplyWithOwner provides a mock function with given fields: doc, namespace, resource
func (_m *mockK8sClient) ApplyWithOwner(doc apply.YamlDocument, namespace string, resource v1.Object) error {
	ret := _m.Called(doc, namespace, resource)

	var r0 error
	if rf, ok := ret.Get(0).(func(apply.YamlDocument, string, v1.Object) error); ok {
		r0 = rf(doc, namespace, resource)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockK8sClient_ApplyWithOwner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyWithOwner'
type mockK8sClient_ApplyWithOwner_Call struct {
	*mock.Call
}

// ApplyWithOwner is a helper method to define mock.On call
//  - doc apply.YamlDocument
//  - namespace string
//  - resource v1.Object
func (_e *mockK8sClient_Expecter) ApplyWithOwner(doc interface{}, namespace interface{}, resource interface{}) *mockK8sClient_ApplyWithOwner_Call {
	return &mockK8sClient_ApplyWithOwner_Call{Call: _e.mock.On("ApplyWithOwner", doc, namespace, resource)}
}

func (_c *mockK8sClient_ApplyWithOwner_Call) Run(run func(doc apply.YamlDocument, namespace string, resource v1.Object)) *mockK8sClient_ApplyWithOwner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(apply.YamlDocument), args[1].(string), args[2].(v1.Object))
	})
	return _c
}

func (_c *mockK8sClient_ApplyWithOwner_Call) Return(_a0 error) *mockK8sClient_ApplyWithOwner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockK8sClient_ApplyWithOwner_Call) RunAndReturn(run func(apply.YamlDocument, string, v1.Object) error) *mockK8sClient_ApplyWithOwner_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockK8sClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockK8sClient creates a new instance of mockK8sClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockK8sClient(t mockConstructorTestingTnewMockK8sClient) *mockK8sClient {
	mock := &mockK8sClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
