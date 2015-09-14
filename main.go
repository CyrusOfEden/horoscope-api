package main

import (
	"github.com/gin-gonic/gin"
    "github.com/robfig/cron"
    "github.com/PuerkitoBio/goquery"
)

import (
	"encoding/json"
	"fmt"
	"log"
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
	doc, _ := goquery.NewDocument(url)
	return parsePage(doc)
}

func seed() [][]horoscope {
	doc, err := goquery.NewDocument(onionHoroscopeUrl)
	if err != nil {
		log.Fatal(err)
	}

	hs := [][]horoscope{}
	doc.Find(".reading-list-item").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("data-absolute-url")
		hs = append(hs, fetchHoroscopes("http://"+url))
	})
	return hs
}

func main() {
	cr := cron.New()
	cr.AddFunc("@weekly", func() {})
	cr.Start()

	app := gin.Default()

	app.GET("/current", func(c *gin.Context) {
		y, w := time.Now().ISOWeek()
		c.Data(200, "application/json", wc[y][w])
	})
	app.GET("/current/:sign", func(c *gin.Context) {
		y, w := time.Now().ISOWeek()
		c.Data(200, "application/json", sc[y][w][c.Param("sign")])
	})
	app.GET("/archive/:year/:week", func(c *gin.Context) {
		y, _ := strconv.Atoi(c.Param("year"))
		w, _ := strconv.Atoi(c.Param("week"))
		c.Data(200, "application/json", wc[y][w])
	})
	app.GET("/archive/:year/:week/:sign", func(c *gin.Context) {
		y, _ := strconv.Atoi(c.Param("year"))
		w, _ := strconv.Atoi(c.Param("week"))
		c.Data(200, "application/json", sc[y][w][c.Param("sign")])
	})

	app.Run(":8000")
}
