package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHealthResponse contains the data that is send to the client requesting the health GET endpoint.
type GetHealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type requestHandler struct {
	AppVersion string `json:"app_version"`
}

func newRequestHandler(appVersion string) requestHandler {
	return requestHandler{AppVersion: appVersion}
}

// getHealth responses with a message indicating the health status of the application. This endpoint can be used to
// determine if the application is reachable.
func (r *requestHandler) getHealth(c *gin.Context) {
	response := &GetHealthResponse{
		Status:  "healthy",
		Version: r.AppVersion,
	}

	c.JSON(http.StatusOK, response)
}
