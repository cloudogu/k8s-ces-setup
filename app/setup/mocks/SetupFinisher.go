// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SetupFinisher is an autogenerated mock type for the SetupFinisher type
type SetupFinisher struct {
	mock.Mock
}

// FinishSetup provides a mock function with given fields:
func (_m *SetupFinisher) FinishSetup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
