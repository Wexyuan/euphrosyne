package response

import (
	"net/http"

	"github.com/Wexyuan/kairos/pkg/errs"
)

// Response is the unified response structure.
type Response struct {
	Code    int    `json:"code"`           // success code or business error code
	Message string `json:"message"`        // result message
	Data    any    `json:"data,omitempty"` // response payload
}

func OK(data any) Response {
	return Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
}

func Fail(err error) Response {
	if businessErr, ok := errs.From(err); ok {
		return Response{
			Code:    businessErr.Code,
			Message: businessErr.Message,
		}
	}
	return Response{
		Code:    http.StatusInternalServerError,
		Message: "server internal error",
	}
}
