package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

func periodicallyUpdate() *time.Ticker {
	ticker := time.NewTicker(time.Hour * 24 * 7)
	go func() {
		for t := range ticker.C {
			FetchHoroscopes()
			data, err := t.MarshalText()
			check(err)
			ioutil.WriteFile("last_update.txt", data, 0644)
		}
	}()
	return ticker
}

func main() {
	seedPtr := flag.Bool("seed", false, "perform the initial seed")
	serverPtr := flag.Bool("server", false, "start the web server")
	portPtr := flag.String("port", "8080", "port to run the web server")

	flag.Parse()

	if *seedPtr {
		fmt.Println("Fetching all horoscopes...")
		SeedHoroscopes()
	} else {
		fmt.Println("Fetching current horoscopes...")
		FetchHoroscopes()
	}

	fmt.Println("Done!")

	if *serverPtr {
		fmt.Println("Scheduling periodic updates... ")
		periodicallyUpdate()

		fmt.Println("Bootstrapping server...")
		Server().Run(":" + *portPtr)
	}
}
