package main

import (
	"flag"
	"fmt"
)

func main() {
	seedPtr := flag.Bool("seed", false, "perform the initial seed")
	pathPtr := flag.String("path", "data", "output path")

	flag.Parse()

	if *seedPtr == true {
		fmt.Println("Fetching all horoscopes...")
		SeedHoroscopes(*pathPtr)
	} else {
		fmt.Println("Fetching current horoscopes...")
		FetchHoroscopes(*pathPtr)
	}

	fmt.Println("Done!")
}
