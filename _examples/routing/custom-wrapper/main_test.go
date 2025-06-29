package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hidevopsio/iris/httptest"
)

type resource string

func (r resource) String() string {
	return string(r)
}

func (r resource) strip(strip string) string {
	s := r.String()
	return strings.TrimPrefix(s, strip)
}

func (r resource) loadFromBase(dir string) string {
	filename := r.String()

	if filename == "/" {
		filename = "/index.html"
	}

	fullpath := filepath.Join(dir, filename)

	b, err := ioutil.ReadFile(fullpath)
	if err != nil {
		panic(fullpath + " failed with error: " + err.Error())
	}

	return string(b)
}

var urls = []resource{
	"/",
	"/index.html",
	"/app.js",
	"/css/main.css",
}

func TestCustomWrapper(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	for _, u := range urls {
		url := u.String()
		contents := u.loadFromBase("./public")

		e.GET(url).Expect().
			Status(httptest.StatusOK).
			Body().Equal(contents)
	}
}
