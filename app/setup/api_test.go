package setup

import (
	"bytes"
	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSetupAPI(t *testing.T) {
	t.Run("should fail because of invalid config paths", func(t *testing.T) {
		// given
		routesMock := newMockGinRoutes(t)
		routesMock.EXPECT().POST("/api/v1/setup", mock.AnythingOfType("gin.HandlerFunc")).RunAndReturn(func(_ string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
			for _, f := range handlerFunc {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				f(c)
			}

			return routesMock
		})
		restConfig := &rest.Config{}
		clientSet := fake.NewSimpleClientset()

		setupCtxBuilder := context.NewSetupContextBuilder("development")
		setupCtxBuilder.DevSetupConfigPath = "invalid"
		t.Setenv("POD_NAMESPACE", "ecosystem")
		t.Setenv("STAGE", "development")

		// when
		logs, err := captureLogs(func() {
			SetupAPI(routesMock, restConfig, clientSet, setupCtxBuilder)
		})

		// then
		require.NoError(t, err)
		assert.Contains(t, logs, "could not find configuration at invalid")
	})
	t.Run("should fail because of failing auth in dogu registry", func(t *testing.T) {
		// given
		routesMock := newMockGinRoutes(t)
		routesMock.EXPECT().POST("/api/v1/setup", mock.AnythingOfType("gin.HandlerFunc")).RunAndReturn(func(_ string, handlerFunc ...gin.HandlerFunc) gin.IRoutes {
			for _, f := range handlerFunc {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				f(c)
			}

			return routesMock
		})
		restConfig := &rest.Config{}
		clientSet := fake.NewSimpleClientset()

		// get project root
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Join(filepath.Dir(b), "../..")

		setupCtxBuilder := context.NewSetupContextBuilder("development")
		setupCtxBuilder.DevSetupConfigPath = filepath.Join(basepath, context.SetupConfigConfigmapDevPath)
		setupCtxBuilder.DevStartupConfigPath = filepath.Join(basepath, context.SetupStartUpConfigMapDevPath)
		setupCtxBuilder.DevDoguRegistrySecretPath = filepath.Join(basepath, "k8s/dev-resources/dogu-registry-secret.example.yaml")
		t.Setenv("POD_NAMESPACE", "ecosystem")
		t.Setenv("STAGE", "development")

		// redirect stdout
		orig := os.Stdout
		defer func() { os.Stdout = orig }()
		r, w, err := os.Pipe()
		require.NoError(t, err)
		defer require.NoError(t, w.Close())
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)

		// when
		logs, err := captureLogs(func() {
			SetupAPI(routesMock, restConfig, clientSet, setupCtxBuilder)
		})

		// then
		require.NoError(t, err)
		assert.Contains(t, logs, "failed to register dogu installation steps: failed to generate dogu step generator: failed to get latest version of dogu [official/ldap]")
	})
}

func captureLogs(f func()) (string, error) {
	realOut := logrus.StandardLogger().Out
	defer logrus.SetOutput(realOut)

	fakeReaderPipe, fakeWriterPipe, err := routeOutputToReplacement()
	if err != nil {
		return "", err
	}

	f()

	actualLogOutput, err := readOutput(fakeReaderPipe, fakeWriterPipe)
	if err != nil {
		return "", err
	}

	return actualLogOutput, nil
}

func routeOutputToReplacement() (readerPipe, writerPipe *os.File, err error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	logrus.SetOutput(w)
	return r, w, nil
}

func readOutput(fakeReaderPipe, fakeWriterPipe *os.File) (string, error) {
	outC := make(chan string)
	errC := make(chan error)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, fakeReaderPipe)
		if err != nil {
			errC <- err
		} else {
			outC <- buf.String()
		}
	}()

	// back to normal state
	err := fakeWriterPipe.Close()
	if err != nil {
		return "", err
	}

	select {
	case actualOutput := <-outC:
		return actualOutput, nil
	case err = <-errC:
		return "", err
	}
}
