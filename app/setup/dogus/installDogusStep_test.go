package dogus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/cesapp-lib/core"
	v1 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
)

var testCtx = context.Background()

func TestNewInstallDogusStep(t *testing.T) {
	t.Run("create new dogu install step", func(t *testing.T) {
		// given
		ecoSystemClientMock := newMockEcoSystemClient(t)
		myDogu := &core.Dogu{Name: "MyName"}

		// when
		installStep := NewInstallDogusStep(ecoSystemClientMock, myDogu, "namespace")

		// then
		require.NotNil(t, installStep)
	})
}

func Test_installDogusStep_GetStepDescription(t *testing.T) {
	t.Run("create new dogu install step", func(t *testing.T) {
		// given
		ecoSystemClientMock := newMockEcoSystemClient(t)
		myDogu := &core.Dogu{Name: "MyName"}
		installStep := NewInstallDogusStep(ecoSystemClientMock, myDogu, "namespace")

		// when
		description := installStep.GetStepDescription()

		// then
		assert.Equal(t, "Installing dogu [MyName]", description)
	})
}

func Test_installDogusStep_PerformSetupStep(t *testing.T) {
	t.Run("failed to get version", func(t *testing.T) {
		// given
		ecoSystemClientMock := newMockEcoSystemClient(t)
		myDogu := &core.Dogu{Name: "MyName", Version: "-----------"}
		installStep := NewInstallDogusStep(ecoSystemClientMock, myDogu, "namespace")

		// when
		err := installStep.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get version from dogu")
	})

	t.Run("should fail on creating dogu", func(t *testing.T) {
		// given
		doguCr := &v1.Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "MyName",
				Namespace: "namespace",
				Labels: map[string]string{
					"app":       "ces",
					"dogu.name": "MyName",
				},
			},
			Spec: v1.DoguSpec{
				Name:    "MyName",
				Version: "1.1.1-1",
			},
			Status: v1.DoguStatus{},
		}

		doguClientMock := NewDoguInterface(t)
		doguClientMock.EXPECT().Create(context.Background(), doguCr, metav1.CreateOptions{}).Return(nil, assert.AnError)
		ecoSystemClientMock := newMockEcoSystemClient(t)
		ecoSystemClientMock.EXPECT().Dogus("namespace").Return(doguClientMock)

		myDogu := &core.Dogu{Name: "MyName", Version: "1.1.1-1"}
		installStep := NewInstallDogusStep(ecoSystemClientMock, myDogu, "namespace")

		// when
		err := installStep.PerformSetupStep(testCtx)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to apply dogu MyName")
	})

	t.Run("successfully apply dogu cr", func(t *testing.T) {
		// given
		doguCr := &v1.Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "MyName",
				Namespace: "namespace",
				Labels: map[string]string{
					"app":       "ces",
					"dogu.name": "MyName",
				},
			},
			Spec: v1.DoguSpec{
				Name:    "MyName",
				Version: "1.1.1-1",
			},
			Status: v1.DoguStatus{},
		}

		doguClientMock := NewDoguInterface(t)
		doguClientMock.EXPECT().Create(context.Background(), doguCr, metav1.CreateOptions{}).Return(doguCr, nil)
		ecoSystemClientMock := newMockEcoSystemClient(t)
		ecoSystemClientMock.EXPECT().Dogus("namespace").Return(doguClientMock)

		myDogu := &core.Dogu{Name: "MyName", Version: "1.1.1-1"}
		installStep := NewInstallDogusStep(ecoSystemClientMock, myDogu, "namespace")

		// when
		err := installStep.PerformSetupStep(testCtx)

		// then
		require.NoError(t, err)
	})
}
