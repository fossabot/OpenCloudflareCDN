package response

import (
	"net/http"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func New(msg string, data ...any) *Response {
	switch len(data) {
	case 0:
		return &Response{
			Msg: msg,
		}
	case 1:
		return &Response{
			Msg:  msg,
			Data: data[0],
		}
	default:
		return &Response{
			Msg:  msg,
			Data: data,
		}
	}
}

func (r *Response) Write(ctx *gin.Context, status ...int) {
	s := http.StatusOK
	if len(status) > 0 {
		s = status[0]
	}

	if !ctx.Writer.Written() {
		ctx.JSON(s, r)
	} else {
		log.Instance.Warn("Duplicate response write attempt",
			zap.String("ctx", util.GinContextString(ctx)))
	}

	ctx.Abort()
}
