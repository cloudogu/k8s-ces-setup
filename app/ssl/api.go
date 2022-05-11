package ssl

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/cesapp-lib/registry"

	"github.com/cloudogu/k8s-ces-setup/app/context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const endpointPostGenerateSSL = "/api/v1/ssl"

// SetupAPI setups the REST API for ssl generation
func SetupAPI(router gin.IRoutes, setupContext *context.SetupContext) {
	logrus.Debugf("Register endpoint [%s][%s]", http.MethodPost, endpointPostGenerateSSL)

	router.POST(endpointPostGenerateSSL, func(ctx *gin.Context) {
		validDays := ctx.Query("days")
		i, err := strconv.ParseInt(validDays, 10, 0)
		if err != nil {
			handleError(ctx, http.StatusBadRequest, err, "Expire days can't convert to integer")
			return
		}

		etcdRegistry, err := registry.New(core.Registry{
			Type:      "etcd",
			Endpoints: []string{fmt.Sprintf("http://%s:4001", setupContext.AppConfig.TargetNamespace)},
		})
		if err != nil {
			handleError(ctx, http.StatusBadRequest, err, "Creating etcd registry")
			return
		}

		config := etcdRegistry.GlobalConfig()
		fqdn, err := config.Get("fqdn")
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, err, "Failed to get FQDN from global config")
			return
		}

		domain, err := config.Get("domain")
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, err, "Failed to get DOMAIN from global config")
			return
		}

		sslGenerator := NewSSLGenerator()
		cert, key, err := sslGenerator.GenerateSelfSignedCert(fqdn, domain, int(i))
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, err, "Failed to generate self-signed certificate and key")
			return
		}

		sslWriter := NewSSLWriter(config)
		err = sslWriter.WriteCertificate("selfsigned", cert, key)
		if err != nil {
			handleError(ctx, http.StatusInternalServerError, err, "Failed to write certificate to global config")
			return
		}

		ctx.Status(http.StatusOK)
	})
}

func handleError(ginCtx *gin.Context, httpCode int, err error, causingAction string) {
	logrus.Error(err.Error())
	ginCtx.String(httpCode, "HTTP %d: An error occurred during this action: %s",
		httpCode, causingAction)
	ginCtx.Writer.WriteHeaderNow()
	ginCtx.Abort()
	_ = ginCtx.Error(err)
}
