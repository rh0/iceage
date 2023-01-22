package main

import (
	"fmt"
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

	mastoURL := os.Getenv("MASTODON_API")
	fmt.Println(mastoURL)

	c := toot.NewClient(os.Getenv("ACCESS_TOKEN"))
	if err = c.PostToot(); err != nil {
		log.Fatal("There was a problem tooting")
	}
}
