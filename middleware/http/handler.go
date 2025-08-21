package http

import (
	"net"
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Handler(httpsPort string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request

		hostOnly, _, err := net.SplitHostPort(req.Host)
		if err != nil {
			hostOnly = req.Host
		}

		targetHost := hostOnly
		if httpsPort != "" && httpsPort != "443" {
			targetHost = net.JoinHostPort(hostOnly, httpsPort)
		}

		u := *req.URL
		u.Scheme = "https"
		u.Host = targetHost

		target := u.String()
		log.Instance.Info("Redirecting to HTTPS", zap.String("url", target), zap.String("ctx", util.GinContextString(ctx)))
		ctx.Header("Location", target)
		response.New(
			"follow url to https",
			gin.H{
				"url": target,
			},
		).Write(ctx, http.StatusPermanentRedirect)
	}
}
