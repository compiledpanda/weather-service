package openweathermap

type Client struct {
}

func (c *Client) GetWeather(lat float64, lon float64) {

}

func NewClient(apiKey string) *Client {
	return &Client{}
}
