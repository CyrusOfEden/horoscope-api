package main

import (
	"flag"
	"fmt"
)

func main() {
	seedPtr := flag.Bool("seed", false, "perform the initial seed")

	flag.Parse()

	if *seedPtr == true {
		fmt.Println("Fetching all horoscopes...")
		SeedHoroscopes()
	} else {
		fmt.Println("Fetching current horoscopes...")
		FetchHoroscopes()
	}

	fmt.Println("Done!")
}
