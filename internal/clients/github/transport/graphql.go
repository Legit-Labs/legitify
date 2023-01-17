package transport

import (
	"net/http"
)

const (
	experimentalApiAcceptHeader = "application/vnd.github.hawkgirl-preview+json"
)

type graphQL struct {
	Base         http.RoundTripper
	acceptHeader string
}

func (t graphQL) RoundTrip(request *http.Request) (*http.Response, error) {
	req2 := CloneRequest(*request)

	req2.Header.Set("Accept", t.acceptHeader)

	return t.Base.RoundTrip(&req2)
}

func CloneRequest(req http.Request) http.Request {
	req.Header = CloneHeader(req.Header)

	return req
}

func CloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}

func NewGraphQL(base http.RoundTripper) *http.Client {
	return &http.Client{Transport: &graphQL{
		acceptHeader: experimentalApiAcceptHeader,
		Base:         base,
	}}
}
