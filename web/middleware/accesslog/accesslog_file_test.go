package accesslog

import (
	"GinLearning/web"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Log_file(t *testing.T) {
	// 创建一个Builder对象
	builder := NewBuilder()

	// log 文件路径
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}

	defer logFile.Close()

	mdl := builder.LogFunc(func(log string) {
		logFile.WriteString(fmt.Sprintf("%s: %s \n", time.Now().Format(time.RFC3339), log))
	}).Build()

	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/hello", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("Hello log!"))
	})

	server.Start(":8081")
}
