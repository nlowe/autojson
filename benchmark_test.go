package autojson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type exampleType struct {
	Foo string
}

func exampleHandler(_ HeaderProvider, r *http.Request) *exampleType {
	return &exampleType{Foo: r.RequestURI}
}

func BenchmarkStdlib(b *testing.B) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(exampleHandler(w, r)); err != nil {
			panic(err)
		}
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		w := &httptest.ResponseRecorder{}
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		h(w, r)
	}
}

func BenchmarkJson(b *testing.B) {
	h := HandlerFunc(exampleHandler)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		w := &httptest.ResponseRecorder{}
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		h(w, r)
	}
}
