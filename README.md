# autojson

[![Build Status](https://travis-ci.org/nlowe/autojson.svg?branch=master)](https://travis-ci.org/nlowe/autojson) [![Coverage Status](https://coveralls.io/repos/github/nlowe/autojson/badge.svg?branch=master)](https://coveralls.io/github/nlowe/autojson?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/nlowe/autojson)](https://goreportcard.com/report/github.com/nlowe/autojson) [![Documentation](https://godoc.org/github.com/nlowe/autojson?status.svg)](https://godoc.org/github.com/nlowe/autojson)

`autojson` makes it easier to return JSON responses and errors for APIs.

## Usage

To add `autojson` to your project, simply `go get` it:

```bash
go get -u github.com/nlowe/autojson
```

`autojson` accepts a wide variety of handler functions. See the [godoc](https://godoc.org/github.com/nlowe/autojson#HandlerFunc)
for details. In general, your handlers now return concrete types and an optional HTTP Status Code and/or error,
`autojson` takes care of the serialization. See the [example](./sample/main.go) for more details.

## Building

This project makes use of Go Modules. You need Go 1.12+ to fully utilize it.

## License

`autojson` is licensed under the MIT License. It is inspired by the AWS Lambda Go SDK, which is licensed under the
Apache 2 License and can be found at https://github.com/aws/aws-lambda-go.
