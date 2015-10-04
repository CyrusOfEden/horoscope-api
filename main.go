package main

import (
	"flag"
	"fmt"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func OutputPath() string {
	return "horoscopes"
}

func periodically(fn func(time.Time)) *time.Ticker {
	ticker := time.NewTicker(time.Hour * 24 * 7)
	go func() {
		for t := range ticker.C {
			fn(t)
		}
	}()
	return ticker
}

func main() {
	seedPtr := flag.Bool("seed", false, "perform the initial seed")
	serverPtr := flag.Bool("server", false, "start the web server")
	portPtr := flag.String("port", ":8000", "port to run the web server")

	flag.Parse()

	if *seedPtr {
		fmt.Println("Fetching all horoscopes...")
		SeedHoroscopes()
	} else {
		fmt.Println("Fetching current horoscopes...")
		FetchHoroscopes()
	}

	if *serverPtr {
		fmt.Print("Setting up store...")
		s := Store()
		s.BuildIndexes()
		fmt.Println(" done!")

		fmt.Print("Scheduling periodic updates... ")
		periodically(func(t time.Time) {
			FetchHoroscopes()
			s.BuildIndexes()
		})
		fmt.Println(" done!")

		fmt.Println("Bootstrapping server...")
		Server(s).Run(*portPtr)
	}
}
