package nws

import (
	"fmt"
	"net/http"
)

type Forecast struct {
	Properties ForecastDetails `json:"properties"`
}

type ForecastDetails struct {
	Updated string           `json:"updated"`
	Periods []ForecastPeriod `json:"periods"`
}

type ForecastPeriod struct {
	Name             string `json:"name"`
	DetailedForecast string `json:"detailedForecast"`
}

func (c *Client) FetchForecast() (Forecast, error) {
    var res Forecast
    getURL := fmt.Sprintf("%s/zones/forecast/%s/forecast", c.BaseURL, c.CountyCode)

    req, err := http.NewRequest("GET", getURL, nil)
    if err != nil {
        return res, err
    }

    if err = c.sendRequest(req, &res); err != nil {
        return res, err
    }

    return res, nil
}

func (c *Client) FormatForecast(forecastResponce Forecast, periods int) string {
    var forecastText string

    for d:=0; d<periods; d++ {
        forecastText += fmt.Sprintf("-= %s =-\n%s\n\n", forecastResponce.Properties.Periods[d].Name, forecastResponce.Properties.Periods[d].DetailedForecast)
    }

    return forecastText
}
