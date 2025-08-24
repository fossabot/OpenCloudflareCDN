package handler

import (
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func APIHealth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Instance.Info("H >> Health", zap.String("ctx", util.GinContextString(ctx)))

		response.New("ok", gin.H{
			"ts": time.Now().UTC().Format(time.RFC3339Nano),
		}).Write(ctx)
	}
}
