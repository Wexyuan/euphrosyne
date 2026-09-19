package base

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Wexyuan/euphrosyne/internal/constant"
	"github.com/Wexyuan/euphrosyne/internal/ecode"
	"github.com/Wexyuan/euphrosyne/pkg/errs"
	"github.com/Wexyuan/euphrosyne/pkg/logger"
	"github.com/Wexyuan/euphrosyne/pkg/response"
)

// Handler provides the shared handler helpers.
type Handler struct {
	log *logger.Logger // shared logger
}

func NewHandler(log *logger.Logger) *Handler {
	return &Handler{log: log}
}

// OK writes a success response.
func (h *Handler) OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response.OK(data))
}

// Fail writes a failure response.
func (h *Handler) Fail(c *gin.Context, httpCode int, err error) {
	if _, ok := errs.From(err); !ok {
		h.log.Errorw("handle request error",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"error", err,
		)
	}
	c.JSON(httpCode, response.Fail(err))
}

// ParsePage parses the pagination query.
func (h *Handler) ParsePage(c *gin.Context) PageQuery {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	return PageQuery{
		Page:     page,
		PageSize: size,
	}.clamp()
}

// ParseID parses the ID path parameter.
func (h *Handler) ParseID(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		h.Fail(c, http.StatusBadRequest, ecode.ErrBadRequest)
		return 0, false
	}
	return id, true
}

// BindJSON binds and validates the request body.
func (h *Handler) BindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		h.Fail(c, http.StatusBadRequest, ecode.ErrValidation)
		return false
	}
	return true
}

// GetUserID returns the authenticated user ID.
func (h *Handler) GetUserID(c *gin.Context) (int64, bool) {
	value, exists := c.Get(string(constant.AuthUserIDKey))
	if !exists {
		return 0, false
	}

	id, ok := value.(int64)
	return id, ok && id > 0
}

// Logger returns the logger.
func (h *Handler) Logger() *logger.Logger {
	return h.log
}
