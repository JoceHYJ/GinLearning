package recovery

import (
	"GinLearning/web"
	"fmt"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "panic 发生了",
		LogFunc: func(ctx *web.Context) {
			fmt.Printf("panic 路径: %s", ctx.Req.URL.String())
		},
	}
	s := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	s.Get("/user", func(ctx *web.Context) {
		panic("panic 发生了")
	})
	s.Start(":8081")
}
