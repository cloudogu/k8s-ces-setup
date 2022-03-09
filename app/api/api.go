package api

import (
	"fmt"

	"github.com/cloudogu/k8s-ces-setup/app/config"
	"github.com/cloudogu/k8s-ces-setup/app/health"
	"github.com/cloudogu/k8s-ces-setup/app/setup"
	"github.com/gin-gonic/gin"
)

// SetupAPI configures the individual endpoints of the API
func SetupAPI(router gin.IRoutes, appConfig config.Config) error {
	err := health.SetupAPI(router, appConfig)
	if err != nil {
		return fmt.Errorf("failed to setup 'health' api; %w", err)
	}

	err = setup.SetupAPI(router, appConfig)
	if err != nil {
		return fmt.Errorf("failed to setup 'setup' api; %w", err)
	}

	return nil
}
