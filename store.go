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

func (s *store) getHoroscopesIndex(y, w string) (p string, k string) {
	p = path.Join(y, w, "index.json")
	k = encodeKey(y, w)
	return
}

func (s *store) getHoroscopeIndex(y, w, sign string) (p string, k string) {
	p = path.Join(y, w, sign+".json")
	k = encodeKey(y, w, sign)
	return
}

func (s *store) getAndUpdate(filePath, key string) ([]byte, bool) {
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
	p, k := s.getHoroscopesIndex(y, w)
	return s.getAndUpdate(p, k)
}

func (s *store) GetHoroscope(y, w, sign string) ([]byte, bool) {
	p, k := s.getHoroscopeIndex(y, w, sign)
	return s.getAndUpdate(p, k)
}

func Store() *store {
	return &store{Cache: cache.New(time.Hour, 5*time.Minute)}
}
