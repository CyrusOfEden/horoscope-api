package main

import (
	"flag"
	"fmt"
)

func OutputPath() string {
	return "horoscopes"
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

	fmt.Println(" done!")

	if *serverPtr {
		fmt.Println("Bootstrapping server...")
		s := Server()
		s.Run(":" + *portPtr)
	}
}
