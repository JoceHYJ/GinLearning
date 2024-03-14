package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个 before")
				next(ctx)
				fmt.Println("第一个 after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个 before")
				next(ctx)
				fmt.Println("第二个 after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个, 看不到这里")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
