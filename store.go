package main

import (
	"github.com/pmylund/go-cache"
)

import (
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type store struct {
	Cache     *cache.Cache
	yearIndex map[string]string
	weekIndex map[string]map[string]string
}

func (s *store) encodeKey(parts ...string) string {
	return strings.Join(parts, "/")
}

func (s *store) getPath(y, w, filename string) string {
	return path.Join(y, w, filename+".json")
}

func (s *store) BuildIndexes() {
	var cwd, yn, ysp, wn, wsp, prevn string
	var ws []int
	var n, w, prev, curr, wi int
	var err error
	var _ys, _ws []os.FileInfo

	cwd, err = os.Getwd()
	check(err)

	ysp = path.Join(cwd, OutputPath())

	_ys, err = ioutil.ReadDir(ysp)
	check(err)

	s.yearIndex = make(map[string]string, len(_ys))
	s.weekIndex = make(map[string]map[string]string, len(_ys))

	for _, y := range _ys {
		yn = y.Name()
		wsp = path.Join(ysp, yn)
		_ws, err = ioutil.ReadDir(wsp)
		check(err)

		ws = make([]int, len(_ws))
		for i, w := range _ws {
			wi, _ = strconv.Atoi(w.Name())
			ws[i] = wi
		}
		sort.Ints(ws)

		s.yearIndex[yn] = yn
		s.weekIndex[yn] = make(map[string]string, len(ws))

		prev = ws[0]
		curr = 0

		for _, w = range ws {
			wn = strconv.Itoa(w)
			s.weekIndex[yn][wn] = wn

			curr, _ = strconv.Atoi(wn)
			if curr - prev >= 1 {
				prevn = strconv.Itoa(prev)
				for n = prev; n <= curr; n++ {
					s.weekIndex[yn][strconv.Itoa(n)] = prevn
				}
			}
			prev = curr
		}

		prevn = strconv.Itoa(prev)
		for n = prev; n <= 52; n++ {
			s.weekIndex[yn][strconv.Itoa(n)] = prevn
		}
	}
}

func (s *store) getIndexes(y, w string) (string, string) {
	y, _ = s.yearIndex[y]
	w, _ = s.weekIndex[y][w]
	return y, w
}

func (s *store) getHoroscopesIndex(y, w string) (string, string) {
	y, w = s.getIndexes(y, w)
	return s.getPath(y, w, "index"), s.encodeKey(y, w)
}

func (s *store) getHoroscopeIndex(y, w, sign string) (string, string) {
	y, w = s.getIndexes(y, w)
	return s.getPath(y, w, sign), s.encodeKey(y, w, sign)
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
