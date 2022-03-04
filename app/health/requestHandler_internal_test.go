package health

import (
	"encoding/json"
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testConfig = config.Config{
	LogLevel: logrus.DebugLevel,
}

func setupGinAndAPIWithPermissions(t *testing.T, config config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	err := SetupAPI(r, config)
	require.NoError(t, err)
	return r
}

func Test_newRequestHandler(t *testing.T) {
	handler := newRequestHandler()

	require.NotNil(t, handler)
}

func Test_requestHandler_getHealth(t *testing.T) {
	// given
	r := setupGinAndAPIWithPermissions(t, testConfig)

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", EndpointGetHealth, nil)
	r.ServeHTTP(w, req)

	// then
	require.Equal(t, http.StatusOK, w.Code)
	var response GetHealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
}
