package main

import (
	"github.com/pmylund/go-cache"
)

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func encodeKey(parts ...string) string {
	return strings.Join(parts, ":")
}

type store struct {
	Cache *cache.Cache
}

func (s *store) getAndUpdate(key string, filePath string) ([]byte, bool) {
	if data, found := s.Cache.Get(key); found {
		return data.([]byte), true
	}

	cwd, err := os.Getwd()
	check(err)
	filePath = path.Join(cwd, OutputPath(), filePath)

	if data, err := ioutil.ReadFile(filePath); err == nil {
		s.Cache.Set(key, data, cache.DefaultExpiration)
		return data, true
	} else {
		return nil, false
	}
}

func (s *store) GetHoroscopes(y, w string) ([]byte, bool) {
	return s.getAndUpdate(encodeKey(y, w), path.Join(y, w, "index.json"))
}

func (s *store) GetHoroscope(y, w, sign string) ([]byte, bool) {
	return s.getAndUpdate(encodeKey(y, w, sign), path.Join(y, w, sign+".json"))
}

func Store() *store {
	return &store{Cache: cache.New(time.Hour, 5*time.Minute)}
}
