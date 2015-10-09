package main

import (
	"github.com/gin-gonic/gin"
)

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func getISOWeek(c *gin.Context) (string, string) {
	return c.MustGet("year").(string), c.MustGet("week").(string)
}

func bootstrapMiddleware(s *store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set("store", s)
		c.Set("human", c.DefaultQuery("human", "true") == "true")
		c.Next()
	}
}

func latestMiddleware(c *gin.Context) {
	y, w := time.Now().ISOWeek()
	c.Set("year", strconv.Itoa(y))
	c.Set("week", strconv.Itoa(w))
	c.Next()
}

func paramsMiddleware(c *gin.Context) {
	c.Set("year", c.Param("year"))
	c.Set("week", c.Param("week"))
	c.Next()
}

func humanMiddleware(c *gin.Context) {
	s := c.MustGet("store").(*store)
	h := c.MustGet("human").(bool)
	y, w := getISOWeek(c)
	if !h && !s.HasWeekStrict(y, w) {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.Next()
	}
}

func doc(filename string) func(*gin.Context) {
	cwd, err := os.Getwd()
	check(err)
	doc, err := ioutil.ReadFile(path.Join(cwd, filename+".html"))
	check(err)
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", doc)
	}
}

func mount(r *gin.RouterGroup) {
	r.GET("", func(c *gin.Context) {
		s := c.MustGet("store").(*store)
		y, w := getISOWeek(c)
		if data, found := s.GetHoroscopes(y, w); found {
			c.Data(http.StatusOK, "application/json", data)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	r.GET("/:sign", func(c *gin.Context) {
		s := c.MustGet("store").(*store)
		y, w := getISOWeek(c)
		if data, found := s.GetHoroscope(y, w, c.Param("sign")); found {
			c.Data(http.StatusOK, "application/json", data)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})
}

func Server(s *store, release bool) *gin.Engine {
	if release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/help", doc("help"))

	l := r.Group("/current")
	l.Use(bootstrapMiddleware(s))
	l.Use(latestMiddleware)
	l.Use(humanMiddleware)
	mount(l)

	a := r.Group("/archive/:year/:week")
	a.Use(bootstrapMiddleware(s))
	a.Use(paramsMiddleware)
	a.Use(humanMiddleware)
	mount(a)

	return r
}
