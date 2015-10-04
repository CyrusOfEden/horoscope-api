package main

import (
	"encoding/json"
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
	releasePtr := flag.Bool("release", false, "set the mode to release")
	debugPtr := flag.Bool("debug", false, "enable debugging mode")

	flag.Parse()

	seed := *seedPtr
	server := *serverPtr
	port := *portPtr
	release := *releasePtr
	debug := *debugPtr

	if seed {
		fmt.Println("Fetching all horoscopes...")
		SeedHoroscopes()
	} else {
		fmt.Println("Fetching current horoscopes...")
		FetchHoroscopes()
	}

	if server {
		fmt.Print("Setting up store... ")
		s := Store()
		s.BuildIndexes()
		fmt.Println("done!")

		if debug {
			data, _ := json.MarshalIndent(s.weekIndex, "", "  ")
			fmt.Println(string(data), "", "  ")
		}

		fmt.Print("Scheduling periodic updates... ")
		periodically(func(t time.Time) {
			FetchHoroscopes()
			s.BuildIndexes()
		})
		fmt.Println("done!")

		fmt.Println("Bootstrapping server... done!")
		if debug {
			Server(s, false).Run(port)
		} else {
			Server(s, release).Run(port)
		}
	}
}
