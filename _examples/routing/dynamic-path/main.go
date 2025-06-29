package main

import (
	"regexp"

	"github.com/hidevopsio/iris"
)

func main() {
	app := iris.New()

	// At the previous example "routing/basic",
	// we've seen static routes, group of routes, subdomains, wildcard subdomains, a small example of parameterized path
	// with a single known paramete and custom http errors, now it's time to see wildcard parameters and macros.

	// Iris, like net/http std package registers route's handlers
	// by a Handler, the iris' type of handler is just a func(ctx iris.Context)
	// where context comes from github.com/hidevopsio/iris/context.
	//
	// Iris has the easiest and the most powerful routing process you have ever meet.
	//
	// At the same time,
	// Iris has its own interpeter(yes like a programming language)
	// for route's path syntax and their dynamic path parameters parsing and evaluation,
	// We call them "macros" for shortcut.
	// How? It calculates its needs and if not any special regexp needed then it just
	// registers the route with the low-level underline  path syntax,
	// otherwise it pre-compiles the regexp and adds the necessary middleware(s).
	//
	// Standard macro types for parameters:
	//  +------------------------+
	//  | {param:string}         |
	//  +------------------------+
	// string type
	// anything
	//
	//  +-------------------------------+
	//  | {param:int} or {param:int} |
	//  +-------------------------------+
	// int type
	// both positive and negative numbers, any number of digits (ctx.Params().GetInt will limit the digits based on the host arch)
	//
	// +-------------------------------+
	// | {param:int64} or {param:long} |
	// +-------------------------------+
	// int64 type
	// -9223372036854775808 to 9223372036854775807
	//
	// +------------------------+
	// | {param:uint8}          |
	// +------------------------+
	// uint8 type
	// 0 to 255
	//
	//
	// +------------------------+
	// | {param:uint64}         |
	// +------------------------+
	// uint64 type
	// 0 to 18446744073709551615
	//
	// +---------------------------------+
	// | {param:bool} or {param:boolean} |
	// +---------------------------------+
	// bool type
	// only "1" or "t" or "T" or "TRUE" or "true" or "True"
	// or "0" or "f" or "F" or "FALSE" or "false" or "False"
	//
	//  +------------------------+
	//  | {param:alphabetical}   |
	//  +------------------------+
	// alphabetical/letter type
	// letters only (upper or lowercase)
	//
	//  +------------------------+
	//  | {param:file}           |
	//  +------------------------+
	// file type
	// letters (upper or lowercase)
	// numbers (0-9)
	// underscore (_)
	// dash (-)
	// point (.)
	// no spaces ! or other character
	//
	//  +------------------------+
	//  | {param:path}           |
	//  +------------------------+
	// path type
	// anything, should be the last part, more than one path segment,
	// i.e: /path1/path2/path3 , ctx.Params().Get("param") == "/path1/path2/path3"
	//
	// if type is missing then parameter's type is defaulted to string, so
	// {param} == {param:string}.
	//
	// If a function not found on that type then the `string` macro type's functions are being used.
	//
	//
	// Besides the fact that iris provides the basic types and some default "macro funcs"
	// you are able to register your own too!.
	//
	// Register a named path parameter function:
	// app.Macros().Number.RegisterFunc("min", func(argument int) func(paramValue string) bool {
	//  [...]
	//  return true/false -> true means valid.
	// })
	//
	// at the func(argument ...) you can have any standard type, it will be validated before the server starts
	// so don't care about performance here, the only thing it runs at serve time is the returning func(paramValue string) bool.
	//
	// {param:string equal(iris)} , "iris" will be the argument here:
	// app.Macros().String.RegisterFunc("equal", func(argument string) func(paramValue string) bool {
	// 	return func(paramValue string) bool { return argument == paramValue }
	// })

	// you can use the "string" type which is valid for a single path parameter that can be anything.
	app.Get("/username/{name}", func(ctx iris.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("name"))
	}) // type is missing = {name:string}

	// Let's register our first macro attached to uint64 macro type.
	// "min" = the function
	// "minValue" = the argument of the function
	// func(uint64) bool = our func's evaluator, this executes in serve time when
	// a user requests a path which contains the :uint64 macro parameter type with the min(...) macro parameter function.
	app.Macros().Get("uint64").RegisterFunc("min", func(minValue uint64) func(uint64) bool {
		// type of "paramValue" should match the type of the internal macro's evaluator function, which in this case is "uint64".
		return func(paramValue uint64) bool {
			return paramValue >= minValue
		}
	})

	// http://localhost:8080/profile/id>=20
	// this will throw 404 even if it's found as route on : /profile/0, /profile/blabla, /profile/-1
	// macro parameter functions are optional of course.
	app.Get("/profile/{id:uint64 min(20)}", func(ctx iris.Context) {
		// second parameter is the error but it will always nil because we use macros,
		// the validaton already happened.
		id := ctx.Params().GetUint64Default("id", 0)
		ctx.Writef("Hello id: %d", id)
	})

	// to change the error code per route's macro evaluator:
	app.Get("/profile/{id:uint64 min(1)}/friends/{friendid:uint64 min(1) else 504}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		friendid := ctx.Params().GetUint64Default("friendid", 0)
		ctx.Writef("Hello id: %d looking for friend id: ", id, friendid)
	}) // this will throw e 504 error code instead of 404 if all route's macros not passed.

	// :uint8 0 to 255.
	app.Get("/ages/{age:uint8 else 400}", func(ctx iris.Context) {
		age, _ := ctx.Params().GetUint8("age")
		ctx.Writef("age selected: %d", age)
	})

	// Another example using a custom regexp or any custom logic.

	// Register your custom argument-less macro function to the :string param type.
	latLonExpr := "^-?[0-9]{1,3}(?:\\.[0-9]{1,10})?$"
	latLonRegex, err := regexp.Compile(latLonExpr)
	if err != nil {
		panic(err)
	}

	// MatchString is a type of func(string) bool, so we use it as it is.
	app.Macros().Get("string").RegisterFunc("coordinate", latLonRegex.MatchString)

	app.Get("/coordinates/{lat:string coordinate() else 502}/{lon:string coordinate() else 502}", func(ctx iris.Context) {
		ctx.Writef("Lat: %s | Lon: %s", ctx.Params().Get("lat"), ctx.Params().Get("lon"))
	})

	//

	// Another one is by using a custom body.
	app.Macros().Get("string").RegisterFunc("range", func(minLength, maxLength int) func(string) bool {
		return func(paramValue string) bool {
			return len(paramValue) >= minLength && len(paramValue) <= maxLength
		}
	})

	app.Get("/limitchar/{name:string range(1,200)}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef(`Hello %s | the name should be between 1 and 200 characters length
		otherwise this handler will not be executed`, name)
	})

	//

	// Register your custom macro function which accepts a slice of strings `[...,...]`.
	app.Macros().Get("string").RegisterFunc("has", func(validNames []string) func(string) bool {
		return func(paramValue string) bool {
			for _, validName := range validNames {
				if validName == paramValue {
					return true
				}
			}

			return false
		}
	})

	app.Get("/static_validation/{name:string has([kataras,gerasimos,maropoulos]}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef(`Hello %s | the name should be "kataras" or "gerasimos" or "maropoulos"
		otherwise this handler will not be executed`, name)
	})

	//

	// http://localhost:8080/game/a-zA-Z/level/42
	// remember, alphabetical is lowercase or uppercase letters only.
	app.Get("/game/{name:alphabetical}/level/{level:int}", func(ctx iris.Context) {
		ctx.Writef("name: %s | level: %s", ctx.Params().Get("name"), ctx.Params().Get("level"))
	})

	app.Get("/lowercase/static", func(ctx iris.Context) {
		ctx.Writef("static and dynamic paths are not conflicted anymore!")
	})

	// let's use a trivial custom regexp that validates a single path parameter
	// which its value is only lowercase letters.

	// http://localhost:8080/lowercase/anylowercase
	app.Get("/lowercase/{name:string regexp(^[a-z]+)}", func(ctx iris.Context) {
		ctx.Writef("name should be only lowercase, otherwise this handler will never executed: %s", ctx.Params().Get("name"))
	})

	// http://localhost:8080/single_file/app.js
	app.Get("/single_file/{myfile:file}", func(ctx iris.Context) {
		ctx.Writef("file type validates if the parameter value has a form of a file name, got: %s", ctx.Params().Get("myfile"))
	})

	// http://localhost:8080/myfiles/any/directory/here/
	// this is the only macro type that accepts any number of path segments.
	app.Get("/myfiles/{directory:path}", func(ctx iris.Context) {
		ctx.Writef("path type accepts any number of path segments, path after /myfiles/ is: %s", ctx.Params().Get("directory"))
	}) // for wildcard path (any number of path segments) without validation you can use:
	// /myfiles/*

	// "{param}"'s performance is exactly the same of ":param"'s.

	// alternatives -> ":param" for single path parameter and "*" for wildcard path parameter.
	// Note these:
	// if  "/mypath/*" then the parameter name is "*".
	// if  "/mypath/{myparam:path}" then the parameter has two names, one is the "*" and the other is the user-defined "myparam".

	// WARNING:
	// A path parameter name should contain only alphabetical letters or digits. Symbols like  '_' are NOT allowed.
	// Last, do not confuse `ctx.Params()` with `ctx.Values()`.
	// Path parameter's values can be retrieved from `ctx.Params()`,
	// context's local storage that can be used to communicate between handlers and middleware(s) can be stored to `ctx.Values()`.
	app.Run(iris.Addr(":8080"))
}
