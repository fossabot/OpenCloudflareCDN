package errorhandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestError_ResponseAlreadyWritten(t *testing.T) {
	t.Parallel()

	log.Instance = zap.NewNop()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})

	testErr := errors.New("test error after response")

	util.GINError(ctx, testErr)

	Error()(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.NotContains(t, w.Body.String(), "oops, something went wrong")
}

func TestError_MultipleErrors(t *testing.T) {
	t.Parallel()

	log.Instance = zap.NewNop()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	testErr1 := errors.New("first error")
	testErr2 := errors.New("second error")

	util.GINError(ctx, testErr1)
	util.GINError(ctx, testErr2)

	Error()(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "oops, something went wrong")
	assert.Contains(t, w.Body.String(), "traceID")
	assert.Len(t, w.Result().Header.Get("Content-Type"), 1)
}

func TestError_NoErrors(t *testing.T) {
	t.Parallel()

	log.Instance = zap.NewNop()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	Error()(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}
