package autojson

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerFunc(t *testing.T) {
	v := "foo"
	tests := []struct {
		name string
		f    interface{}

		status   int
		response string
	}{
		{name: "struct", f: func() *string { return &v }, status: http.StatusOK, response: `"foo"`},
		{name: "struct_with_status", f: func() (int, *string) { return http.StatusTeapot, &v }, status: http.StatusTeapot, response: `"foo"`},
		{name: "struct_with_error", f: func() (*string, error) { return &v, nil }, status: http.StatusOK, response: `"foo"`},
		{name: "struct_with_status_and_error", f: func() (int, *string, error) { return http.StatusCreated, &v, nil }, status: http.StatusCreated, response: `"foo"`},
		{name: "error", f: func() (*string, error) { return nil, fmt.Errorf("bar") }, status: http.StatusInternalServerError, response: `{"code":500,"error":"bar"}`},
		{name: "error_with_status", f: func() (int, *string, error) { return http.StatusBadRequest, nil, fmt.Errorf("baz") }, status: http.StatusBadRequest, response: `{"code":400,"error":"baz"}`},

		{name: "in_struct", f: func(_ HeaderProvider, r *http.Request) *string { return &r.RequestURI }, status: http.StatusOK, response: `"/"`},
		{name: "in_struct_with_status", f: func(_ HeaderProvider, r *http.Request) (int, *string) { return http.StatusTeapot, &r.RequestURI }, status: http.StatusTeapot, response: `"/"`},
		{name: "in_struct_with_error", f: func(_ HeaderProvider, r *http.Request) (*string, error) { return &r.RequestURI, nil }, status: http.StatusOK, response: `"/"`},
		{name: "in_struct_with_status_and_error", f: func(_ HeaderProvider, r *http.Request) (int, *string, error) {
			return http.StatusCreated, &r.RequestURI, nil
		}, status: http.StatusCreated, response: `"/"`},
		{name: "in_error", f: func(_ HeaderProvider, r *http.Request) (*string, error) {
			return nil, fmt.Errorf("%sbar", r.RequestURI)
		}, status: http.StatusInternalServerError, response: `{"code":500,"error":"/bar"}`},
		{name: "in_error_with_status", f: func(_ HeaderProvider, r *http.Request) (int, *string, error) {
			return http.StatusBadRequest, nil, fmt.Errorf("%sbaz", r.RequestURI)
		}, status: http.StatusBadRequest, response: `{"code":400,"error":"/baz"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h := HandlerFunc(tt.f)

			h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

			assert.Equal(t, tt.status, w.Code)

			body, _ := ioutil.ReadAll(w.Body)
			assert.Regexp(t, tt.response, string(body))
		})
	}
}

func TestHandler_CanUseCustomErrorConverter(t *testing.T) {
	w := httptest.NewRecorder()
	h := HandlerFunc(func() (*string, error) { return nil, fmt.Errorf("asdf") })

	defer UseDefaultErrorConverter()
	UseErrorConverter(func(_ int, _ *http.Request, _ error) interface{} {
		return &struct {
			Foo string `json:"foo"`
		}{Foo: "bar"}
	})

	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	body, _ := ioutil.ReadAll(w.Body)
	assert.Regexp(t, `{"foo":"bar"}`, string(body))
}

func TestHandlerFunc_PanicsForBadHandler(t *testing.T) {
	require.Panics(t, func() {
		_ = HandlerFunc("foo")
	})
}

type deadWriter struct{}

func (d *deadWriter) Header() http.Header {
	return http.Header{}
}

func (d *deadWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("dummy error")
}

func (d *deadWriter) WriteHeader(statusCode int) {}

func TestHandlerFunc_PanicsIfJsonEncodeFails(t *testing.T) {
	v := ""
	h := HandlerFunc(func() *string { return &v })

	require.Panics(t, func() {
		h.ServeHTTP(&deadWriter{}, httptest.NewRequest(http.MethodGet, "/", nil))
	})
}
