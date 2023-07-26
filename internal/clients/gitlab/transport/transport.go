package transport

import (
	"crypto/tls"
	"net/http"

	"github.com/Legit-Labs/legitify/internal/clients/transport"
)

func NewHttpClient(ignoreInvalidCertificate bool) *http.Client {
	if ignoreInvalidCertificate {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &http.Client{
		Transport: transport.NewCacheTransport(),
	}
}
