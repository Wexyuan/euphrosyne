package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/Wexyuan/euphrosyne/internal/ecode"
	"github.com/Wexyuan/euphrosyne/pkg/logger"
	"github.com/Wexyuan/euphrosyne/pkg/response"
)

// Recovery creates the panic recovery middleware.
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Attach the stack as a field so a multi-line panic stays greppable.
				log.Errorw("recover from panic", "panic", r, "stack", string(debug.Stack()))
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Fail(ecode.ErrInternalServer))
			}
		}()

		c.Next()
	}
}
