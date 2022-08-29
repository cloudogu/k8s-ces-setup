// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ResourceRegistryClient is an autogenerated mock type for the resourceRegistryClient type
type ResourceRegistryClient struct {
	mock.Mock
}

// GetResourceFileContent provides a mock function with given fields: resourceURL
func (_m *ResourceRegistryClient) GetResourceFileContent(resourceURL string) ([]byte, error) {
	ret := _m.Called(resourceURL)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(resourceURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(resourceURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}