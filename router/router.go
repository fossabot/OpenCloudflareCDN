package router

import (
	"github.com/Sn0wo2/OpenCloudflareCDN/router/errorhandler"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/handler"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/notfound"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/proxy"
	"github.com/Sn0wo2/OpenCloudflareCDN/router/static"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".mp4", ".mp3", ".zip", ".gz", ".br"})))
	router.Use(errorhandler.Error())

	api := router.Group("/ocfc-api")

	v0 := api.Group("/v0")
	{
		v0.Any("/error", handler.APIError())
	}

	v1 := api.Group("/v1")
	{
		v1.Any("/health", handler.APIHealth())
		v1.Any("/info", handler.APIInfo())
		v1.Any("/verify", handler.APIVerify())
	}

	router.NoRoute(func(ctx *gin.Context) {
		proxy.Proxy()(ctx)

		if ctx.Writer.Written() {
			return
		}

		static.Handle()(ctx)

		if ctx.Writer.Written() {
			return
		}

		notfound.Handler()(ctx)
	})
}
