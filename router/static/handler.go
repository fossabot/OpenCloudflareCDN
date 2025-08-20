package static

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.Next()
			return
		}

		staticRoot := config.Instance.StaticPath
		if staticRoot == "" {
			_ = ctx.Error(errors.New("static path is empty"))
			ctx.Next()
			return
		}

		status := http.StatusOK

		if absStaticRoot, err := filepath.Abs(staticRoot); err == nil {
			if absTargetPath, err := filepath.Abs(filepath.Join(staticRoot, strings.TrimPrefix(ctx.Request.URL.Path, "/"))); err == nil {
				if absStaticRoot == absTargetPath {
					status = http.StatusForbidden
				}
			}
		}

		filePath := GetStaticFile(config.Instance.StaticIndex, staticRoot, ctx.Request.URL.Path)
		if filePath == "" {
			_ = ctx.Error(errors.New("static file not found"))
			ctx.Next()
			return
		}

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			_ = ctx.Error(err)
			ctx.Next()
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		if strings.HasPrefix(contentType, "text/") {
			contentType += "; charset=utf-8"
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			ctx.Next()
			return
		}

		log.Instance.Info("S >> Static file served", zap.String("file", filePath), zap.String("ctx", util.GinContextString(ctx)))

		ctx.Data(status, contentType, fileContent)
	}
}
