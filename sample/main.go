package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nlowe/autojson"
)

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/struct", autojson.HandlerFunc(exampleStruct))
	m.HandleFunc("/structcode", autojson.HandlerFunc(exampleCustomResponse))
	m.HandleFunc("/error", autojson.HandlerFunc(exampleError))
	m.HandleFunc("/errorcode", autojson.HandlerFunc(exampleCustomError))

	s := http.Server{Handler: m, Addr: "0.0.0.0:5000"}

	go func() {
		fmt.Println("Starting Up...")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Shutting Down...")
	if err := s.Shutdown(ctx); err != nil {
		panic(err)
	}
}

type exampleResponse struct {
	Foo  string              `json:"foo"`
	Bar  int                 `json:"bar"`
	Fizz map[string][]string `json:"headers"`
}

func exampleStruct(w autojson.HeaderProvider, r *http.Request) *exampleResponse {
	w.Header().Add("X-My-Custom-Header", "autojson/sample")
	return &exampleResponse{
		Foo:  "Hello, world!",
		Bar:  rand.Int(),
		Fizz: r.Header,
	}
}

func exampleCustomResponse(w autojson.HeaderProvider, r *http.Request) (int, *exampleResponse) {
	return http.StatusTeapot, exampleStruct(w, r)
}

func exampleError() (*exampleResponse, error) {
	return nil, fmt.Errorf("something went wrong")
}

func exampleCustomError() (int, *exampleResponse, error) {
	return http.StatusBadRequest, nil, fmt.Errorf("something went wrong (with a custom response code)")
}
