package server

import (
	"net/http"

	"github.com/compiledpanda/weatherservice/internal/config"
	"github.com/compiledpanda/weatherservice/internal/endpoint"
	"github.com/compiledpanda/weatherservice/internal/openweathermap"
)

func New(c config.Config) *http.ServeMux {
	// Create clients
	owm := openweathermap.NewClient(c.OpenWeatherMapKey)

	// Create mux
	mux := http.NewServeMux()

	// Add endpoints
	// Our weather endpoint
	mux.HandleFunc("/v1/conditions", endpoint.NotAllowed())
	mux.HandleFunc("GET /v1/conditions", endpoint.GetConditions(owm))

	// General 404 handler for everything else
	mux.HandleFunc("/", endpoint.NotFound())

	return mux
}
