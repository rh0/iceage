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

	// Get Weather Alerts!
	alerts := nws.AlertList{}
	w := nws.NewClient()
	alerts, err = w.FetchAlerts()
	if err != nil {
		log.Fatal("There was a problem fetching alerts from NWS: ", err)
	}

	//Try to toot!
	t := toot.NewClient(os.Getenv("ACCESS_TOKEN"))
	if err = t.PostToot(alerts.Features[0].Properties.Headline, alerts.Features[0].Properties.Description); err != nil {
		log.Fatal("There was a problem tooting: ", err)
	}
}
