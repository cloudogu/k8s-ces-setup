// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SSLWriter is an autogenerated mock type for the SSLWriter type
type SSLWriter struct {
	mock.Mock
}

// WriteCertificate provides a mock function with given fields: certType, cert, key
func (_m *SSLWriter) WriteCertificate(certType string, cert string, key string) error {
	ret := _m.Called(certType, cert, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(certType, cert, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
