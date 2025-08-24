package util

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGINError(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		err         error
		expectError bool
	}{
		{
			name:        "should add error to context and abort",
			err:         errors.New("test error"),
			expectError: true,
		},
		{
			name:        "should handle nil error gracefully",
			err:         nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			GINError(ctx, tt.err)

			if tt.expectError {
				if len(ctx.Errors) == 0 {
					t.Error("Expected error to be added to context, but no errors found")
				} else if !errors.Is(ctx.Errors[0].Err, tt.err) || len(ctx.Errors) != 1 {
					t.Errorf("Expected error %v, got %v", tt.err, ctx.Errors)
				}
			} else if len(ctx.Errors) != 0 {
				t.Errorf("Expected no errors for nil error, but got %d errors", len(ctx.Errors))
			}
		})
	}
}

func TestGINError_Integration(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("should work in handler context", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		testErr := errors.New("integration test error")

		handler := func(ctx *gin.Context) {
			err := testErr
			if err != nil {
				GINError(ctx, err)

				return
			}

			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		}

		handler(ctx)

		if len(ctx.Errors) == 0 {
			t.Error("Expected error to be added to context")
		} else if !errors.Is(testErr, ctx.Errors[0].Err) || len(ctx.Errors) != 1 {
			t.Errorf("Expected error %v, got %v", testErr, ctx.Errors)
		}

		if w.Body.Len() > 0 {
			t.Errorf("Expected empty response body, but got: %s", w.Body.String())
		}
	})
}

func TestGINError_MultipleErrors(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	t.Run("should handle multiple errors correctly", func(t *testing.T) {
		t.Parallel()

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		err1 := errors.New("first error")
		err2 := errors.New("second error")

		GINError(ctx, err1)

		if len(ctx.Errors) != 1 {
			t.Errorf("Expected 1 error after first GINError call, got %d", len(ctx.Errors))
		}

		GINError(ctx, err2)

		if len(ctx.Errors) != 2 {
			t.Errorf("Expected 2 errors after second GINError call, got %d", len(ctx.Errors))
		}

		if !errors.Is(err1, ctx.Errors[0].Err) {
			t.Errorf("Expected first error to be %v, got %v", err1, ctx.Errors[0].Err)
		}

		if !errors.Is(err2, ctx.Errors[1].Err) {
			t.Errorf("Expected second error to be %v, got %v", err2, ctx.Errors[1].Err)
		}
	})
}
