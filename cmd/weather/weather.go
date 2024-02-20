package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/compiledpanda/weatherservice/internal/config"
	"github.com/compiledpanda/weatherservice/internal/server"
)

func main() {
	// Set log level to debug for now
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Load Config
	// TODO this would ideally be wired up to a true config service and pull from
	// something like a configMap, SSM, a config file, and/or env vars
	c := config.Config{
		Port:              8080,
		OpenWeatherMapKey: os.Getenv("OPEN_WEATHER_MAP_KEY"),
	}

	// Validate config
	// TODO this would ideally be handled in a true config service
	if c.OpenWeatherMapKey == "" {
		slog.Error("OPEN_WEATHER_MAP_KEY env variable must be set")
		os.Exit(1)
	}

	// Create Server
	// TODO normally we would create a signal/context to pass through a interrupts for a graceful
	// shutdown, but the specifics depend on how it is run (traditional, docker, lambda, etc...)
	s := server.New(c)

	// Start Server
	// TODO we could go https here, but again that depends on the context in which this will be deployed
	slog.Info("Starting Server", "port", c.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", c.Port), s)
}
