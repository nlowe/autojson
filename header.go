package autojson

import "net/http"

// HeaderProvider is a stripped down http.ResponseWriter that only provides access
// to the headers to be sent in the response
type HeaderProvider interface {
	Header() http.Header
}
