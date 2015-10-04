package main

import (
	"github.com/gin-gonic/gin"
)

import (
	"net/http"
	"strconv"
	"time"
)

func getDate(c *gin.Context) (string, string) {
	return c.MustGet("year").(string), c.MustGet("week").(string)
}

func mountRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		s := c.MustGet("store").(*store)
		year, week := getDate(c)
		if horoscopes, found := s.GetHoroscopes(year, week); found {
			c.Data(http.StatusOK, "application/json", horoscopes)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})
	r.GET("/:sign", func(c *gin.Context) {
		s := c.MustGet("store").(*store)
		year, week := getDate(c)
		sign := c.Param("sign")
		if horoscope, found := s.GetHoroscope(year, week, sign); found {
			c.Data(http.StatusOK, "application/json", horoscope)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})
}

func Server(s *store) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("store", s)
		c.Next()
	})

	l := r.Group("/latest")
	l.Use(func(c *gin.Context) {
		y, w := time.Now().ISOWeek()
		c.Set("year", strconv.Itoa(y))
		c.Set("week", strconv.Itoa(w))
		c.Next()
	})
	mountRoutes(l)

	a := r.Group("/archive/:year/:week")
	a.Use(func(c *gin.Context) {
		c.Set("year", c.Param("year"))
		c.Set("week", c.Param("week"))
		c.Next()
	})
	mountRoutes(a)

	return r
}
