package main

import (
	"github.com/PuerkitoBio/goquery"
)

import (
	"encoding/json"
	"fmt"
	"flag"
	"log"
	"strings"
	"strconv"
	"time"
	"path"
	"os"
	"io/ioutil"
)

const dateFormat string = "Jan 2, 2006"
const onionHoroscopeUrl string = "http://www.theonion.com/features/horoscope"

type horoscope struct {
	Year       int    `json:"year"`
	Week       int    `json:"week"`
	Sign       string `json:"sign"`
	Prediction string `json:"prediction"`
}

func parseHoroscope(text string, date time.Time) horoscope {
	ls := strings.Split(strings.Trim(text, " \n"), "\n")
	l1, p := ls[0], strings.Trim(ls[1], " \n")
	ts := strings.Split(l1, " | ")
	s := strings.ToLower(ts[0])
	y, w := date.ISOWeek()
	return horoscope{Sign: s, Prediction: p, Year: y, Week: w}
}

func parseDate(text string) (time.Time, error) {
	t := strings.Trim(text, " \n")
	cidx := strings.Index(t, ",")
	c := t[0:3] + " " + t[cidx-1:]
	return time.Parse(dateFormat, c)
}

func parsePage(doc *goquery.Document) []horoscope {
	date, err := parseDate(doc.Find(".content-published").Text())
	if err == nil {
		fmt.Printf("Fetched horocopes for %s\n", date)
	} else {
		log.Fatal(err)
	}

	hs := make([]horoscope, 12)
	doc.Find(".astro .large-thing").Each(func(i int, s *goquery.Selection) {
		hs[i] = parseHoroscope(s.Text(), date)
	})
	return hs
}

func fetchHoroscopes(url string) []horoscope {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return parsePage(doc)
}

func seedHoroscopes() []horoscope {
	doc, err := goquery.NewDocument(onionHoroscopeUrl)
	if err != nil {
		log.Fatal(err)
	}

	hs := make([]horoscope, 0)
	doc.Find(".reading-list-item").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("data-absolute-url")
		hs = append(hs, fetchHoroscopes("http://"+url)...)
	})
	return hs
}

func cacheHoroscopes(hs []horoscope) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory: ", err)
	}

	for _, h := range hs {
		p := path.Join(cwd, "data", strconv.Itoa(h.Year), strconv.Itoa(h.Week))

		if err := os.MkdirAll(p, 0777); err != nil {
			log.Fatal("Failed to make directory: ", err)
		}

		mrshl, _ := json.Marshal(h)
		ioutil.WriteFile(path.Join(p, h.Sign + ".json"), mrshl, 0644)
	}
}

func main() {
	seedPtr := flag.Bool("seed", false, "perform the initial seed")
	flag.Parse()

	var hs []horoscope
	if *seedPtr == true {
		fmt.Println("Fetching all horoscopes...")
		hs = seedHoroscopes()
	} else {
		fmt.Println("Fetching current horoscopes...")
		hs = fetchHoroscopes(onionHoroscopeUrl)
	}

	fmt.Println("Caching horoscopes...")
	cacheHoroscopes(hs)
	fmt.Println("Done!")
}
