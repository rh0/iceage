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
	alerts := nws.AlertList{}
	alerts, err = w.FetchAlerts()
	if err != nil {
		log.Fatal("There was a problem fetching alerts from NWS: ", err)
	}

    // Get forecast!
    forecast := nws.Forecast{}
    forecast, err = w.FetchForecast()
    if err != nil {
        log.Fatal("There was a problem fetching the forecast from NWS: ", err)
    }


	t := toot.NewClient(os.Getenv("ACCESS_TOKEN"))

	// Try to toot the alert!
	if err = t.PostToot(alerts.Features[0].Properties.Event, alerts.Features[0].Properties.Description); err != nil {
		log.Fatal("There was a problem tooting: ", err)
	}

    // Tey to toot the forecast!
    formattedForecast := w.FormatForecast(forecast)
    if err = t.PostToot("2 day Forecast", formattedForecast); err != nil {
        log.Fatal("There was a problem tooting: ", err)
    }
}
