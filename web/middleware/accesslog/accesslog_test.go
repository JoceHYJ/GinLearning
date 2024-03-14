package accesslog

import (
	"GinLearning/web"
	"fmt"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	// 链式调用
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	//server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello tomato")
	})
	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
