package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewK8sClient(t *testing.T) {
	actual := NewK8sClient(nil)

	require.NotNil(t, actual)
}
