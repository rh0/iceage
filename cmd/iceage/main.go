package main

import (
	"fmt"
	"iceage/pkg/nws"
	//"iceage/pkg/toot"
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

	// Try to toot!
	// t := toot.NewClient(os.Getenv("ACCESS_TOKEN"))
	// if err = t.PostToot("BEEP BOOP!"); err != nil {
	// 	log.Fatal("There was a problem tooting")
	// }

	w := nws.NewClient()
	if err = w.FetchAlert(); err != nil {
		log.Fatal("There was a problem fetching alerts from NWS", err)
	}
}
