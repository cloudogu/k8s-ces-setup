package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockExiter struct {
	Error error `json:"error"`
}

func (e *mockExiter) Exit(err error) {
	e.Error = err
}

func Test_createSetupRouter(t *testing.T) {
	t.Run("Startup without error", func(t *testing.T) {
		// given
		mockExiter := &mockExiter{}

		// when
		createSetupRouter(mockExiter, "testdata/k8s-ces-setup-testdata.yaml")

		//then
		assert.Nil(t, mockExiter.Error)
	})

	t.Run("Startup error", func(t *testing.T) {
		// given
		mockExiter := &mockExiter{}

		// when
		createSetupRouter(mockExiter, "not-a-config")

		//then
		assert.NotNil(t, mockExiter.Error)
		assert.Equal(t, "could not find configuration at not-a-config", mockExiter.Error.Error())
	})
}
