package main

import (
	"github.com/PuerkitoBio/goquery"
)

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const dateFormat string = "Jan 2, 2006"
const onionHoroscopeUrl string = "http://www.theonion.com/features/horoscope"

type horoscope struct {
	Year       int    `json:"year"`
	Week       int    `json:"week"`
	Sign       string `json:"sign"`
	Prediction string `json:"prediction"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
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
	check(err)
	fmt.Printf("Fetched horocopes for %s\n", date)

	hs := make([]horoscope, 12)
	doc.Find(".astro .large-thing").Each(func(i int, s *goquery.Selection) {
		hs[i] = parseHoroscope(s.Text(), date)
	})
	return hs
}

func processUrl(url string) {
	doc, err := goquery.NewDocument(url)
	check(err)
	cacheHoroscopes(parsePage(doc))
}

func cacheHoroscopes(hs []horoscope) {
	cwd, err := os.Getwd()
	check(err)

	p := path.Join(cwd, "data", strconv.Itoa(hs[0].Year), strconv.Itoa(hs[0].Week))
	err = os.MkdirAll(p, 0777)
	check(err)

	data, err := json.Marshal(hs)
	check(err)
	ioutil.WriteFile(path.Join(p, "index.json"), data, 0644)

	for _, h := range hs {
		data, err = json.Marshal(h)
		check(err)
		ioutil.WriteFile(path.Join(p, h.Sign+".json"), data, 0644)
	}
}

// Public API
func FetchHoroscopes() {
	processUrl(onionHoroscopeUrl)
}

func SeedHoroscopes() {
	doc, err := goquery.NewDocument(onionHoroscopeUrl)
	check(err)

	messages := make(chan bool)

	elems := doc.Find(".reading-list-item")

	elems.Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("data-absolute-url")
		go func(url string) {
			processUrl(url)
			messages <- true
		}("http://" + url)
	})

	for i := 0; i < elems.Length(); i++ {
		<-messages
	}
}
