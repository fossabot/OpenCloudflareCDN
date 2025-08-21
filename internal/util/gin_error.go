package util

import (
	"github.com/gin-gonic/gin"
)

// GINError adds an error to the Gin context and aborts the request processing.
// Into a single call, eliminating the need for manual ctx.Next() calls.
//
// Usage:
//
//	if err != nil {
//	    util.GINError(ctx, err)
//	    return
//	}
func GINError(ctx *gin.Context, err error) {
	if err != nil {
		_ = ctx.Error(err)
	}
}
