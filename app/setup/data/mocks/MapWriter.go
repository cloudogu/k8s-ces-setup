// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	context "github.com/cloudogu/k8s-ces-setup/app/context"

	mock "github.com/stretchr/testify/mock"
)

// MapWriter is an autogenerated mock type for the MapWriter type
type MapWriter struct {
	mock.Mock
}

// WriteConfigToMap provides a mock function with given fields: registryConfig
func (_m *MapWriter) WriteConfigToStringDataMap(registryConfig context.CustomKeyValue) (map[string]map[string]string, error) {
	ret := _m.Called(registryConfig)

	var r0 map[string]map[string]string
	if rf, ok := ret.Get(0).(func(context.CustomKeyValue) map[string]map[string]string); ok {
		r0 = rf(registryConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.CustomKeyValue) error); ok {
		r1 = rf(registryConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
