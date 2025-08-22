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
			util.GINError(ctx, errors.New("static path is empty"))

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
			util.GINError(ctx, errors.New("static path not found"))
			ctx.Next()

			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if strings.HasPrefix(contentType, "text/") && !strings.Contains(contentType, "charset") {
			contentType += "; charset=utf-8"
		}

		fileContent, err := os.ReadFile(filePath) //nolint:gosec
		if err != nil {
			util.GINError(ctx, err)
			ctx.Next()

			return
		}

		log.Instance.Info("S >> Static file served", zap.String("file", filePath), zap.String("ctx", util.GinContextString(ctx)))

		ctx.Data(status, contentType, fileContent)
	}
}
