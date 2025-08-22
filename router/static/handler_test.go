package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/OpenCloudflareCDN/config"
	"github.com/Sn0wo2/OpenCloudflareCDN/internal/util"
	"github.com/Sn0wo2/OpenCloudflareCDN/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	if config.Instance == nil {
		config.Instance = &config.Config{
			StaticPath:  "./static",
			StaticIndex: "index.html",
			Log: config.Log{
				Level: "error",
				Dir:   "",
			},
		}
	}

	if log.Instance == nil {
		log.Instance = zap.NewNop()
	}

	os.Exit(m.Run())
}

func setupTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "oops, something went wrong",
				"data": gin.H{
					"traceID": "test-trace-id",
				},
			})
		}
	})
	router.Use(Handle())

	return router
}

func TestStaticHandler_EmptyStaticPath(t *testing.T) {
	t.Parallel()

	originalStaticPath := config.Instance.StaticPath

	defer func() {
		config.Instance.StaticPath = originalStaticPath
	}()

	config.Instance.StaticPath = ""

	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "oops, something went wrong")
}

func TestStaticHandler_FileNotFound(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	originalStaticPath := config.Instance.StaticPath
	originalStaticIndex := config.Instance.StaticIndex

	defer func() {
		config.Instance.StaticPath = originalStaticPath
		config.Instance.StaticIndex = originalStaticIndex
	}()

	config.Instance.StaticPath = tempDir
	config.Instance.StaticIndex = "index.html" //nolint:goconst

	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent.html", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "oops, something went wrong")
}

func TestStaticHandler_FileReadError(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.html")
	require.NoError(t, os.WriteFile(testFile, util.StringToBytes("test content"), 0o000))

	originalStaticPath := config.Instance.StaticPath
	originalStaticIndex := config.Instance.StaticIndex

	defer func() {
		config.Instance.StaticPath = originalStaticPath
		config.Instance.StaticIndex = originalStaticIndex
	}()

	config.Instance.StaticPath = tempDir
	config.Instance.StaticIndex = "index.html"

	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusInternalServerError {
		assert.Contains(t, w.Body.String(), "oops, something went wrong")
	}
}

func TestStaticHandler_SuccessfulFileServing(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	testContent := "<!DOCTYPE html><html><body>Test Content</body></html>"
	testFile := filepath.Join(tempDir, "test.html")
	require.NoError(t, os.WriteFile(testFile, util.StringToBytes(testContent), 0o600))

	originalStaticPath := config.Instance.StaticPath
	originalStaticIndex := config.Instance.StaticIndex

	defer func() {
		config.Instance.StaticPath = originalStaticPath
		config.Instance.StaticIndex = originalStaticIndex
	}()

	config.Instance.StaticPath = tempDir
	config.Instance.StaticIndex = "index.html"

	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/test.html", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testContent, w.Body.String())
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestStaticHandler_NonGetMethod(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	originalStaticPath := config.Instance.StaticPath

	defer func() {
		config.Instance.StaticPath = originalStaticPath
	}()

	config.Instance.StaticPath = tempDir

	router := setupTestRouter()
	req := httptest.NewRequest(http.MethodPost, "/test.html", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}
