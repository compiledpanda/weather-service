.Phony: start
start:
	go run cmd/weather/weather.go

.Phony: test
test:
	go test ./...