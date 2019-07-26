package autojson

import "net/http"

// ErrorConverter defines a function that takes in the response code,
// the request, and the error encountered and returns a custom type
// to be JSON Encoded to customize responses when returning errors
type ErrorConverter func(int, *http.Request, error) interface{}

var (
	defaultErrorConverter = func(status int, _ *http.Request, err error) interface{} {
		return &struct {
			Code    int    `json:"code"`
			Message string `json:"error"`
		}{
			Code:    status,
			Message: err.Error(),
		}
	}

	errorConverter = defaultErrorConverter
)

// Use the provided error converter when returning custom errors.
//
// The default format is:
//  {"code": <response code>, "error": err.Error()}
func UseErrorConverter(c ErrorConverter) {
	errorConverter = c
}

// UseDefaultErrorConverter resets the error converter to the default one
func UseDefaultErrorConverter() {
	errorConverter = defaultErrorConverter
}
