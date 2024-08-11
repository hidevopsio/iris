package main

import (
	"github.com/hidevopsio/iris"

	"github.com/hidevopsio/iris/middleware/pprof"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1> Please click <a href='/debug/pprof'>here</a>")
	})

	app.Any("/debug/pprof/{action:path}", pprof.New())
	//                              ___________
	app.Run(iris.Addr(":8080"))
}
