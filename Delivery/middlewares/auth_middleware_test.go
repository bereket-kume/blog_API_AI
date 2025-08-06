package middlewares

import (
	"blog-api/Domain/models"
	"blog-api/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		mockClaims     *models.UserAccessClaims
		mockError      error
		allowedRoles   []string
		expectedStatus int
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			mockClaims:     nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid header format",
			authHeader:     "Token abc123",
			mockClaims:     nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalidtoken",
			mockClaims:     nil,
			mockError:      errors.New("invalid token"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid token but insufficient role",
			authHeader: "Bearer validtoken",
			mockClaims: &models.UserAccessClaims{
				UserID: "123",
				Email:  "test@example.com",
				Role:   "user",
			},
			mockError:      nil,
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:       "valid token with allowed role",
			authHeader: "Bearer validtoken",
			mockClaims: &models.UserAccessClaims{
				UserID: "123",
				Email:  "test@example.com",
				Role:   "admin",
			},
			mockError:      nil,
			allowedRoles:   []string{"admin", "superadmin"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTokenService := new(mocks.MockTokenService)
			if tt.authHeader != "" && tt.mockError != nil || tt.mockClaims != nil {
				mockTokenService.
					On("VerifyAccessToken", mock.Anything).
					Return(tt.mockClaims, tt.mockError)
			}

			r := gin.New()
			r.Use(AuthMiddleware(mockTokenService, tt.allowedRoles...))
			r.GET("/protected", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
