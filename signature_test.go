package autojson

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Foo struct{}

func TestValidateSignature_Valid(t *testing.T) {
	tests := []struct {
		f interface{}

		in     bool
		status bool
		error  bool
	}{
		{f: func() *Foo { return nil }, in: false, status: false, error: false},
		{f: func() (*Foo, error) { return nil, nil }, in: false, status: false, error: true},
		{f: func() (int, *Foo) { return 0, nil }, in: false, status: true, error: false},
		{f: func() (int, *Foo, error) { return 0, nil, nil }, in: false, status: true, error: true},

		{f: func(HeaderProvider, *http.Request) *Foo { return nil }, in: true, status: false, error: false},
		{f: func(HeaderProvider, *http.Request) (*Foo, error) { return nil, nil }, in: true, status: false, error: true},
		{f: func(HeaderProvider, *http.Request) (int, *Foo) { return 0, nil }, in: true, status: true, error: false},
		{f: func(HeaderProvider, *http.Request) (int, *Foo, error) { return 0, nil, nil }, in: true, status: true, error: true},
	}

	for _, tt := range tests {
		t.Run(dumpSignature(reflect.ValueOf(tt.f).Type()), func(t *testing.T) {
			sig, err := validateSignature(tt.f)

			require.NoError(t, err)
			assert.Equal(t, tt.in, sig.InParameters)
			assert.Equal(t, tt.status, sig.Status)
			assert.Equal(t, tt.error, sig.Error)
		})
	}
}

func TestValidateSignature_Invalid(t *testing.T) {
	tests := []struct {
		name string
		f    interface{}
		err  string
	}{
		{name: "invalid type", f: "asdf", err: "handler kind string is not a func"},
		{name: "not enough outs", f: func() {}, err: "unknown handler signature: func()"},
		{name: "too many outs", f: func() (int, int, int, int) { return 0, 0, 0, 0 }, err: "unknown handler signature: func() (int, int, int, int)"},
		{name: "not enough ins", f: func(int) error { return nil }, err: "unknown handler signature: func(int) (error)"},
		{name: "too many ins", f: func(int, int, int) error { return nil }, err: "unknown handler signature: func(int, int, int) (error)"},
		{name: "in: header: wrong type", f: func(int, *http.Request) error { return nil }, err: "input parameter type mismatch: index 0, got int, want autojson.HeaderHandler [func(int, *http.Request) (error)]"},
		{name: "in: request: wrong type", f: func(HeaderProvider, int) error { return nil }, err: "input parameter type mismatch: index 1, got int, want *http.Request [func(autojson.HeaderProvider, int) (error)]"},
		{name: "out: 2: wrong status type", f: func() (string, string) { return "", "" }, err: "return parameter type mismatch: index 0, got string, want int [func() (string, string)]"},
		{name: "out: 3: wrong status type", f: func() (string, string, error) { return "", "", nil }, err: "return parameter type mismatch: index 0, got string, want int [func() (string, string, error)]"},
		{name: "out: 3: wrong status type", f: func() (int, string, string) { return 0, "", "" }, err: "return parameter type mismatch: index 2, got string, want error [func() (int, string, string)]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateSignature(tt.f)

			require.EqualError(t, err, tt.err)
		})
	}
}
