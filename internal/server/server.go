package server

import (
	"net/http"

	"github.com/compiledpanda/weatherservice/internal/config"
	"github.com/compiledpanda/weatherservice/internal/endpoint"
)

func New(c config.Config) *http.ServeMux {
	// Create clients
	// TODO

	// Create mux
	mux := http.NewServeMux()

	// Add endpoints
	// Our weather endpoint
	mux.HandleFunc("/v1/conditions", endpoint.NotAllowed())
	mux.HandleFunc("GET /conditions", endpoint.GetConditions())

	// General 404 handler for everything else
	mux.HandleFunc("/", endpoint.NotFound())

	return mux
}
