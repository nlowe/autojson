package autojson

import (
	"encoding/json"
	"net/http"
	"reflect"
)

const (
	// The Content-Type Header
	HeaderContentType = "Content-Type"

	// The Content Type for JSON responses
	ContentTypeJSON = "application/json"
)

// HandlerFunc returns a handler wrapper that makes encoding responses from
// handler types easier. It accepts the following signatures:
//
//  * func() TOut
//  * func() (TOut, error)
//  * func() (int, TOut)
//  * func() (int, TOut, error)
//  * func(HeaderProvider, *http.Request) TOut
//  * func(HeaderProvider, *http.Request) (TOut, error)
//  * func(HeaderProvider, *http.Request) (int, TOut)
//  * func(HeaderProvider, *http.Request) (int, TOut, error)
//
// The value of TOut (if specified) must be a pointer or interface type. If
// it is nil, no body is written. If an error is returned, it is encoded
// instead of TOut. If the first return parameter is an int, it is used for
// the HTTP Status code. When omitted, 500 is used for any call that return
// an error, and 200 is used for all other calls.
//
// The handler panics if it fails to encode the JSON response
func HandlerFunc(h interface{}) http.HandlerFunc {
	sig, err := validateSignature(h)
	if err != nil {
		panic(err)
	}

	handler := reflect.ValueOf(h)
	return func(w http.ResponseWriter, r *http.Request) {
		var args []reflect.Value
		if sig.InParameters {
			args = append(args, reflect.ValueOf(w), reflect.ValueOf(r))
		}

		result := handler.Call(args)

		var code int
		if sig.Status {
			code = result[0].Interface().(int)
			result = result[1:]
		} else {
			code = http.StatusOK
		}
		w.Header().Set(HeaderContentType, ContentTypeJSON)

		var response interface{}
		if sig.Error {
			err := result[len(result)-1]

			if !err.IsNil() {
				if !sig.Status {
					code = http.StatusInternalServerError
				}
				response = errorConverter(code, r, err.Interface().(error))
			} else {
				response = result[0].Interface()
			}
		} else if !result[0].IsNil() {
			response = result[0].Interface()
		}

		w.WriteHeader(code)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
}
