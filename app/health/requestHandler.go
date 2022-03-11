package health

import (
	"net/http"

	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/gin-gonic/gin"
)

// GetHealthResponse contains the data that is send to the client requesting the health GET endpoint
type GetHealthResponse struct {
	Status  string         `json:"status"`
	Version config.Version `json:"version"`
}

type requestHandler struct {
}

// getHealth responses with a message indicating the health status of the application. This endpoint can be used to
// determine if the application is reachable.
func (r *requestHandler) getHealth(c *gin.Context) {
	response := &GetHealthResponse{
		Status:  "healthy",
		Version: config.AppVersion,
	}

	c.JSON(http.StatusOK, response)
}
