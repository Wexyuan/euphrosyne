package response

import "github.com/Wexyuan/euphrosyne/pkg/errs"

// Response is the unified response structure.
type Response struct {
	// Code is the success code or the business error code.
	Code int `json:"code"`
	// Message describes the result.
	Message string `json:"message"`
	// Data carries the payload and is omitted when empty.
	Data any `json:"data,omitempty"`
}

// OK builds a success response.
func OK(data any) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// Fail builds a failure response from err.
func Fail(err error) Response {
	code, message := 500, err.Error()
	if businessErr, ok := errs.From(err); ok {
		code, message = businessErr.Code, businessErr.Message
	}
	return Response{
		Code:    code,
		Message: message,
	}
}
