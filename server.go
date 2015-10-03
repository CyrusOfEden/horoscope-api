package main

import (
	"github.com/gin-gonic/gin"
)

import (
	"net/http"
	"strconv"
	"time"
)

func getDate(c *gin.Context) (year, week string) {
	year = c.MustGet("year").(string)
	week = c.MustGet("week").(string)
	return
}

func mountRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		cache := c.MustGet("store").(*store)
		y, w := getDate(c)
		if hs, found := cache.GetHoroscopes(y, w); found {
			c.Data(http.StatusOK, "application/json", hs)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})
	r.GET("/:sign", func(c *gin.Context) {
		cache := c.MustGet("store").(*store)
		y, w := getDate(c)
		s := c.Param("sign")
		if h, found := cache.GetHoroscope(y, w, s); found {
			c.Data(http.StatusOK, "application/json", h)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})
}

func Server() *gin.Engine {
	s := Store()

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
