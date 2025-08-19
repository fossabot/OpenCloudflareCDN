package http

import (
	"net"
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Handler(httpsPort string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		host, _, err := net.SplitHostPort(ctx.Request.Host)
		if err != nil {
			host = ctx.Request.Host
		}

		targetHost := host
		if httpsPort != "" && httpsPort != "443" {
			targetHost = net.JoinHostPort(host, httpsPort)
		}

		targetURL := "https://" + targetHost + ctx.Request.RequestURI
		log.Instance.Info("Redirecting to HTTPS", zap.String("url", targetURL))
		response.New(
			"follow url to https",
			gin.H{
				"url": targetURL,
			},
		).Write(ctx, http.StatusPermanentRedirect)
	}
}
