package endpoint

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/compiledpanda/weatherservice/internal/openweathermap"
)

// TODO I'm definitely not in love with the field names and would want
// to sync with the team/PM to pick better names :)
type GetConditionsResponse struct {
	Temperature float64 `json:"temperature"`
	Units       string  `json:"units"`
	Condition   string  `json:"condition"`
	FeelsLike   string  `json:"feelsLike"`
}

func GetConditions(owm openweathermap.Interface) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Authenticate
		// TODO this will depend on how we protect our API (via Authorization Header, etc...)

		// Authorize
		// TODO this depends on how we determine who can call this endpoints and
		// (potentially) with what parameters

		// Validate Request
		lat, err := getAndValidateLat(req)
		if err != nil {
			returnError(w, http.StatusBadRequest, err.Error())
			return
		}
		lon, err := getAndValidateLon(req)
		if err != nil {
			returnError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Call OpenWeatherMap
		data, err := owm.GetWeather(lat, lon)
		if err != nil {
			returnError(w, http.StatusInternalServerError, "Could Not Get Weather Data")
			return
		}

		// Construct Response
		res := GetConditionsResponse{
			Temperature: data.Main.Temp,
			Units:       "F",
			Condition:   formatConditions(data.Weather),
			FeelsLike:   formatFeelsLike(data.Main.Temp),
		}

		returnJSON(w, http.StatusOK, res)
	}
}

// IMPORTANT: error message is passed to the client and MUST be clean
func getAndValidateLat(req *http.Request) (float64, error) {
	raw := req.URL.Query().Get("lat")
	if raw == "" {
		return 0, errors.New("lat query parameter is required")
	}

	lat, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, errors.New("lat query parameter must be a valid float")
	}

	if lat < -90.0 || lat > 90.0 {
		return 0, errors.New("lat query parameter must be between -90 and 90")
	}

	return lat, nil
}

// IMPORTANT: error message is passed to the client and MUST be clean
func getAndValidateLon(req *http.Request) (float64, error) {
	raw := req.URL.Query().Get("lon")
	if raw == "" {
		return 0, errors.New("lon query parameter is required")
	}

	lon, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, errors.New("lon query parameter must be a valid float")
	}

	if lon < -180.0 || lon > 180.0 {
		return 0, errors.New("lon query parameter must be between -180 and 180")
	}

	return lon, nil
}

func formatConditions(weather []openweathermap.WeatherWeather) string {
	conditions := []string{}
	for _, w := range weather {
		conditions = append(conditions, w.Description)
	}
	return strings.Join(conditions, ", ")
}

func formatFeelsLike(temp float64) string {
	if temp > 80 {
		return "hot"
	}
	if temp < 40 {
		return "cold"
	}
	return "moderate"
}
