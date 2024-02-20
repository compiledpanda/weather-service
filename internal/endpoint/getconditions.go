package endpoint

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/compiledpanda/weatherservice/internal/openweathermap"
)

type GetConditionsResponse struct {
	Temperature int    `json:"temperature"`
	Units       string `json:"units"`
	Condition   string `json:"condition"`
	FeelsLike   string `json:"feelsLike"`
}

func GetConditions(owm *openweathermap.Client) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// Authenticate
		// TODO this will depend on how we protect our API (via Authorization Header, etc...)

		// Authorize
		// TODO this depends on how we determine who can call this endpoints and
		// (potentially) with what parameters

		// Validate Request
		lat, err := getConditionsGetAndValidateLat(req)
		if err != nil {
			returnError(w, http.StatusBadRequest, err.Error())
			return
		}
		lon, err := getConditionsGetAndValidateLon(req)
		if err != nil {
			returnError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Call OpenWeatherMap
		// TODO get data and marshal response
		owm.GetWeather(lat, lon)

		// Construct Response
		// TODO
		res := GetConditionsResponse{}

		returnJSON(w, http.StatusOK, res)
	}
}

// IMPORTANT: error message is passed to the client and MUST be clean
func getConditionsGetAndValidateLat(req *http.Request) (float64, error) {
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
func getConditionsGetAndValidateLon(req *http.Request) (float64, error) {
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
