package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/compiledpanda/weatherservice/internal/openweathermap"
)

type testOWM struct {
	_GetWeather func(lat float64, lon float64) (*openweathermap.WeatherData, error)
}

func (o testOWM) GetWeather(lat float64, lon float64) (*openweathermap.WeatherData, error) {
	return o._GetWeather(lat, lon)
}

func TestGetConditions(t *testing.T) {
	// Mock Client responses
	happyPath := func(lat float64, lon float64) (*openweathermap.WeatherData, error) {
		return &openweathermap.WeatherData{
			Main:    openweathermap.WeatherMain{Temp: 13.1},
			Weather: []openweathermap.WeatherWeather{{Description: "balmy :)"}},
		}, nil
	}
	errorPath := func(lat float64, lon float64) (*openweathermap.WeatherData, error) {
		return nil, errors.New("Boom!")
	}

	tests := []struct {
		test       string
		lat        float64
		lon        float64
		GetWeather func(lat float64, lon float64) (*openweathermap.WeatherData, error)
		status     int
		res        *GetConditionsResponse
	}{
		{"Invalid lat", 200, 0, happyPath, http.StatusBadRequest, nil},
		{"Invalid lon", 0, 200, happyPath, http.StatusBadRequest, nil},
		{"GetWeather Error", 0, 0, errorPath, http.StatusInternalServerError, nil},
		{"Happy Path", 0, 0, happyPath, http.StatusOK, &GetConditionsResponse{
			Temperature: 13.1,
			Units:       "F",
			Condition:   "balmy :)",
			FeelsLike:   "cold",
		}},
	}

	for _, tc := range tests {
		// Setup
		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/conditions?lat=%f&lon=%f", tc.lat, tc.lon), nil)
		w := httptest.NewRecorder()
		handler := GetConditions(testOWM{_GetWeather: tc.GetWeather})

		// Call handler
		handler(w, req)
		res := w.Result()

		// Verify status code
		if res.StatusCode != tc.status {
			t.Errorf("%s: status: expected %d, actual %d", tc.test, tc.status, res.StatusCode)
		}

		// Verify response if status is 200
		if res.StatusCode == http.StatusOK {
			var body GetConditionsResponse
			json.NewDecoder(res.Body).Decode(&body)
			if !reflect.DeepEqual(body, *tc.res) {
				t.Errorf("%s: body: expected %v, actual %v", tc.test, tc.res, body)
			}
		}
	}
}

func TestGetAndValidateLat(t *testing.T) {
	tests := []struct {
		test string
		raw  string
		lat  float64
		err  bool
	}{
		{"missing", "", 0, true},
		{"invalid", "bogus", 0, true},
		{"too large", "91", 0, true},
		{"too small", "-91", 0, true},
		{"valid", "0", 0, false},
		{"valid float", "11.1", 11.1, false},
	}

	for _, tc := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/conditions?lat=%s", tc.raw), nil)
		lat, err := getAndValidateLat(req)
		if tc.err {
			if err == nil {
				t.Errorf("%s: expected error", tc.test)
			}
		} else {
			if lat != tc.lat {
				t.Errorf("%s: expected %f, actual %f", tc.test, tc.lat, lat)
			}
		}
	}
}

func TestGetAndValidateLon(t *testing.T) {
	tests := []struct {
		test string
		raw  string
		lon  float64
		err  bool
	}{
		{"missing", "", 0, true},
		{"invalid", "bogus", 0, true},
		{"too large", "181", 0, true},
		{"too small", "-181", 0, true},
		{"valid", "0", 0, false},
		{"valid float", "11.1", 11.1, false},
	}

	for _, tc := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/conditions?lon=%s", tc.raw), nil)
		lon, err := getAndValidateLon(req)
		if tc.err {
			if err == nil {
				t.Errorf("%s: expected error", tc.test)
			}
		} else {
			if lon != tc.lon {
				t.Errorf("%s: expected %f, actual %f", tc.test, tc.lon, lon)
			}
		}
	}
}

func TestFormatConditions(t *testing.T) {
	tests := []struct {
		test    string
		weather []openweathermap.WeatherWeather
		res     string
	}{
		{"none", nil, ""},
		{"one", []openweathermap.WeatherWeather{
			{Description: "first"},
		}, "first"},
		{"two", []openweathermap.WeatherWeather{
			{Description: "first"},
			{Description: "second"},
		}, "first, second"},
	}

	for _, tc := range tests {
		res := formatConditions(tc.weather)
		if res != tc.res {
			t.Errorf("%s: expected %s, actual %s", tc.test, tc.res, res)
		}
	}
}

func TestFormatFeelsLike(t *testing.T) {
	tests := []struct {
		test string
		temp float64
		res  string
	}{
		{"hot", 80.1, "hot"},
		{"upper moderate", 80, "moderate"},
		{"moderate", 70, "moderate"},
		{"lower moderate", 40, "moderate"},
		{"cold", 39.9, "cold"},
	}

	for _, tc := range tests {
		res := formatFeelsLike(tc.temp)
		if res != tc.res {
			t.Errorf("%s: expected %s, actual %s", tc.test, tc.res, res)
		}
	}
}
