package main

import (
	 "github.com/gin-gonic/gin"
)

import (
	"flag"
	"fmt"
)

func main() {
	fetchPtr := flag.Bool("fetch", true, "fetch files")
	seedPtr := flag.Bool("seed", false, "perform the initial seed")
	serverPtr := flag.Bool("serve", false, "start a static file server")

	flag.Parse()

	if *fetchPtr == true {
		if *seedPtr == true {
			fmt.Println("Fetching all horoscopes...")
			SeedHoroscopes()
		} else {
			fmt.Println("Fetching current horoscopes...")
			FetchHoroscopes()
		}

		fmt.Println("Done!")
	}

	if *serverPtr == true {
		r := gin.Default()
		r.Static("/", "./data")
		r.Run(":8000")
	}
}
