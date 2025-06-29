package main

import (
	"testing"

	"github.com/hidevopsio/iris/httptest"
)

func TestMVCHelloWorld(t *testing.T) {
	e := httptest.New(t, newApp())

	e.GET("/").Expect().Status(httptest.StatusOK).
		ContentType("text/html", "utf-8").Body().Equal("<h1>Welcome</h1>")

	e.GET("/ping").Expect().Status(httptest.StatusOK).
		Body().Equal("pong")

	e.GET("/hello").Expect().Status(httptest.StatusOK).
		JSON().Object().Value("message").Equal("Hello Iris!")

	e.GET("/custom_path").Expect().Status(httptest.StatusOK).
		Body().Equal("hello from the custom handler without following the naming guide")
}
