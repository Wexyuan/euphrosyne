package base

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Wexyuan/euphrosyne/internal/ecode"
	"github.com/Wexyuan/euphrosyne/pkg/logger"
	"github.com/Wexyuan/euphrosyne/pkg/response"
)

// Handler provides the shared helpers for business handlers.
type Handler struct {
	log *logger.Logger
}

// NewHandler creates the base handler.
func NewHandler(log *logger.Logger) *Handler {
	return &Handler{log: log}
}

// OK writes a success response.
func (h *Handler) OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response.OK(data))
}

// Fail writes a failure response.
func (h *Handler) Fail(c *gin.Context, httpCode int, err error) {
	c.JSON(httpCode, response.Fail(err))
}

// ParsePage parse the pagination parameters of the request.
func (h *Handler) ParsePage(c *gin.Context) PageQuery {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	return PageQuery{
		Page:     page,
		PageSize: size,
	}.resolve()
}

// ParseID parse the path parameter as a positive integer.
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

// Logger returns the shared logger.
func (h *Handler) Logger() *logger.Logger {
	return h.log
}
