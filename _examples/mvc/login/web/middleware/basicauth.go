// file: middleware/basicauth.go

package middleware

import "github.com/hidevopsio/iris/middleware/basicauth"

// BasicAuth middleware sample.
var BasicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{
		"admin": "password",
	},
})
