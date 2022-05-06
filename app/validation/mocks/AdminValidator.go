// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	context "github.com/cloudogu/k8s-ces-setup/app/context"
	mock "github.com/stretchr/testify/mock"
)

// AdminValidator is an autogenerated mock type for the AdminValidator type
type AdminValidator struct {
	mock.Mock
}

// ValidateAdmin provides a mock function with given fields: admin, dsType
func (_m *AdminValidator) ValidateAdmin(admin context.User, dsType string) error {
	ret := _m.Called(admin, dsType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.User, string) error); ok {
		r0 = rf(admin, dsType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}