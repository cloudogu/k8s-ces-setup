package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testContext = context.SetupContext{
	AppVersion: "",
	AppConfig: context.Config{
		LogLevel: logrus.DebugLevel,
	},
}

func setupGinAndAPI() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	SetupAPI(r, testContext)
	return r
}

func Test_newRequestHandler(t *testing.T) {
	handler := requestHandler{}

	require.NotNil(t, handler)
}

func Test_requestHandler_getHealth(t *testing.T) {
	// given
	r := setupGinAndAPI()

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", endpointGetHealth, nil)
	r.ServeHTTP(w, req)

	// then
	require.Equal(t, http.StatusOK, w.Code)
	var response GetHealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
}
