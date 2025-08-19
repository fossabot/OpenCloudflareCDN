package static

import (
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet {
			staticRoot := config.Instance.StaticPath
			if staticRoot == "" {
				ctx.Next()

				return
			}

			filePath := GetStaticFile(config.Instance.StaticIndex, staticRoot, ctx.Request.URL.Path)
			if filePath == "" {
				ctx.Next()

				return
			}

			log.Instance.Info("S >> Static file served", zap.String("file", filePath), zap.String("ctx", util.GinContextString(ctx)))
			ctx.File(filePath)
		}
	}
}
