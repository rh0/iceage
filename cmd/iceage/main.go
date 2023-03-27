package main

import (
	"iceage/pkg/nws"
	"iceage/pkg/toot"
	"log"
	"os"
)

func main() {
	// Sanity checks
	mastoURL := os.Getenv("MASTODON_API")
	nwsURL := os.Getenv("NWS_API")
	l := log.Default()
	l.Printf("Using Masto instance: %s", mastoURL)
	l.Printf("Using NWS API: %s", nwsURL)

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
	forecast, err := w.FetchForecast()
	if err != nil {
		l.Fatal("There was a problem fetching the forecast from NWS: ", err)
	}

	l.Println("Forecast Received!")
	t := toot.NewClient(os.Getenv("ACCESS_TOKEN"))

	// Arrange forcast data for tooting
	formattedForecast := w.FormatForecast(forecast, 6)
	forecastToot := toot.Toot{
		Spoiler: "3 day forecast",
		Status:  formattedForecast,
	}

	// Try to toot the forecast!
	if err = t.Toot(forecastToot); err != nil {
		l.Fatal("There was a problem tooting: ", err)
	}
	l.Println("Forecast Tooted!")
}
