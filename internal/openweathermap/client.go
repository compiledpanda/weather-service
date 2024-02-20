package openweathermap

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Key        string
	httpClient *http.Client
}

type WeatherData struct {
	Main    WeatherMain      `json:"main"`
	Weather []WeatherWeather `json:"weather"`
	// Other fields omitted for brevity :)
}

type WeatherMain struct {
	Temp float64 `json:"temp"`
	// Other fields omitted for brevity :)
}

type WeatherWeather struct {
	Description string `json:"description"`
	// Other fields omitted for brevity :)
}

func (c *Client) GetWeather(lat float64, lon float64) (*WeatherData, error) {
	// Marshal query params
	slog.Debug("Getting Weather", "lat", lat, "lon", lon)
	params := url.Values{}
	params.Set("lon", strconv.FormatFloat(lon, 'E', -1, 64))
	params.Set("lat", strconv.FormatFloat(lat, 'E', -1, 64))
	params.Set("units", "imperial")
	params.Set("lang", "en")
	params.Set("appid", c.Key)

	// Call API
	res, err := c.httpClient.Get("https://api.openweathermap.org/data/2.5/weather?" + params.Encode())
	if err != nil {
		slog.Error("Could not get weather", "err", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	// Check Response Code
	if res.StatusCode != http.StatusOK {
		slog.Error("Weather API error", "statusCode", res.StatusCode)
		return nil, fmt.Errorf("weather API returned %d", res.StatusCode)
	}

	// Decode response
	var data WeatherData
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		slog.Error("Could not decode weather response", "err", err.Error())
		return nil, err
	}

	return &data, nil
}

func NewClient(apiKey string) *Client {
	// Ideally this would instantiate a new SDK and place it in the Client struct
	// We are just using the key as extant SDKs appear to be very outdated and we only need a single endpoint
	return &Client{
		Key: apiKey,
		// TODO use an actual SDK or configure this client with production sensible defaults :)
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
