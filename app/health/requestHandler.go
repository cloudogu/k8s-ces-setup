package health

import (
	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/gin-gonic/gin"
)

// GetHealthResponse contains the data that is send to the client requesting the health GET endpoint
type GetHealthResponse struct {
	Status  string         `json:"status"`
	Version config.Version `json:"version"`
}

func newRequestHandler() *requestHandler {
	return &requestHandler{}
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

	c.JSON(200, response)
}
