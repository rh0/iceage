package main

import (
	"fmt"
	"iceage/pkg/nws"
	"iceage/pkg/toot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Sanity checks
	mastoURL := os.Getenv("MASTODON_API")
	nwsURL := os.Getenv("NWS_API")
	fmt.Println(mastoURL)
	fmt.Println(nwsURL)

	w := nws.NewClient()

	// Get Weather Alerts!
	// alerts := nws.AlertList{}
	// alerts, err = w.FetchAlerts()
	// if err != nil {
	// 	log.Fatal("There was a problem fetching alerts from NWS: ", err)
	// }

	// Try to toot the alert!
	// if err = t.PostToot(alerts.Features[0].Properties.Event, alerts.Features[0].Properties.Description); err != nil {
	// 	log.Fatal("There was a problem tooting: ", err)
	// }

	// Get forecast!
	forecast := nws.Forecast{}
	forecast, err = w.FetchForecast()
	if err != nil {
		log.Fatal("There was a problem fetching the forecast from NWS: ", err)
	}

	t := toot.NewClient(os.Getenv("ACCESS_TOKEN"))

	// Arrange forcast data for tooting
	formattedForecast := w.FormatForecast(forecast, 6)
	forecastToot := toot.Toot{
		Spoiler: "3 day forecast",
		Status:  formattedForecast,
	}

	// Try to toot the forecast!
	if err = t.Toot(forecastToot); err != nil {
		log.Fatal("There was a problem tooting: ", err)
	}
}
