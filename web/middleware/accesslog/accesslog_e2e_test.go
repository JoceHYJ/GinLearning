//go:build e2e

package accesslog

import (
	"GinLearning/web"
	"fmt"
	"testing"
)

func TestMiddlewareBuilderE2E(t *testing.T) {
	builder := MiddlewareBuilder{}
	// 链式调用
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	//server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		//fmt.Println("hello tomato"
		ctx.Resp.Write([]byte("hello tomato"))
	})
	server.Start(":8081")
}
