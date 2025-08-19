package handler

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/Sn0wo2/OpenCloudflareCDN/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func APIVerify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			_ = ctx.Error(err)

			return
		}

		res := gjson.Parse(util.BytesToString(bodyBytes))

		log.Instance.Info("V >> Turnstile verification request", zap.String("rayID", res.Get("rayID").String()), zap.String("ctx", util.GinContextString(ctx)))

		resp, err := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", url.Values{"secret": {config.Instance.TurnstileSecretKey}, "response": {res.Get("turnstileToken").String()}})
		if err != nil {
			_ = ctx.Error(err)

			return
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			_ = ctx.Error(err)

			return
		}

		result := gjson.Parse(util.BytesToString(body))
		if !result.Get("success").Bool() {
			log.Instance.Warn(
				"V >> Turnstile verification failed",
				zap.String("rayID", res.Get("rayID").String()),
				zap.String("error-codes", result.Get("error-codes").String()),
				zap.String("ctx", util.GinContextString(ctx)),
			)
			response.New("failed", gin.H{"ec": result.Get("error-codes").String()}).Write(ctx, http.StatusBadRequest)

			return
		}

		age := 24 * time.Hour
		exp := time.Now().Add(age).Unix()

		tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": exp,
			"ip":  ctx.ClientIP(),
			"ua":  ctx.Request.UserAgent(),
		}).SignedString(util.StringToBytes(config.Instance.JWTSecret))
		if err != nil {
			_ = ctx.Error(err)

			return
		}

		ctx.SetCookie("__ocfc_v", tokenStr, int(age.Seconds()), "", "", false, true)

		log.Instance.Info("V >> Turnstile verification success", zap.String("rayID", res.Get("rayID").String()), zap.String("ctx", util.GinContextString(ctx)))

		response.New("success").Write(ctx)
	}
}
