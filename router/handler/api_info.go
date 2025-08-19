package handler

import (
    "crypto/rand"
    "encoding/hex"

    "github.com/Sn0wo2/OpenCloudflareCDN/config"
    "github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
    "github.com/Sn0wo2/OpenCloudflareCDN/log"
    "github.com/Sn0wo2/OpenCloudflareCDN/response"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func APIInfo() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        rayIDBytes := make([]byte, 16)
        _, _ = rand.Read(rayIDBytes)
        rayID := hex.EncodeToString(rayIDBytes)

        log.Instance.Info("I >> Info request", zap.String("rayID", rayID), zap.String("ctx", util.GinContextString(ctx)))

        response.New(
            "success",
            gin.H{
                "rayID":   rayID,
                "siteKey": config.Instance.TurnstileSiteKey,
            },
        ).Write(ctx)
    }
}
