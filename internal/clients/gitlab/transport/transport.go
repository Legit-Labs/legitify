package transport

import (
	"net/http"

	"github.com/Legit-Labs/legitify/internal/clients/transport"
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Transport: transport.NewCacheTransport(),
	}
}
