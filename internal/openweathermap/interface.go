package openweathermap

type Interface interface {
	GetWeather(lat float64, lon float64) (*WeatherData, error)
}
