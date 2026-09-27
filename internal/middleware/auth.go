package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Wexyuan/kairos/internal/constant"
	"github.com/Wexyuan/kairos/internal/ecode"
	"github.com/Wexyuan/kairos/internal/module/auth"
	"github.com/Wexyuan/kairos/pkg/response"
)

// AuthRequired creates the authentication middleware.
func AuthRequired(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		// Keep the scheme check case-insensitive because HTTP auth schemes are not case sensitive.
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ecode.ErrAuthInvalidAccessToken))
			return
		}

		userID, err := tokens.VerifyAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Fail(ecode.ErrAuthInvalidAccessToken))
			return
		}

		c.Set(string(constant.AuthUserIDKey), userID)
		c.Next()
	}
}
