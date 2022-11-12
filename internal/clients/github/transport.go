package github

import "net/http"

type transport struct {
	Base         http.RoundTripper
	AcceptHeader *string
}

func (t transport) RoundTrip(request *http.Request) (*http.Response, error) {
	req2 := CloneRequest(*request)

	if t.AcceptHeader != nil {
		req2.Header.Set("Accept", *t.AcceptHeader)
	}

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

func NewClientWithAcceptHeader(base http.RoundTripper, acceptHeader *string) *http.Client {
	return &http.Client{Transport: &transport{
		AcceptHeader: acceptHeader,
		Base:         base,
	}}
}
