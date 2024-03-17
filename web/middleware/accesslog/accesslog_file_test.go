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
	//builder := NewBuilder()

	// 替换为 logrus 库进行日志记录
	//builder := NewLogrusBuilder()

	// 替换为 zap 库进行日志记录
	builder := NewZapBuilder()

	// log 文件路径
	logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}

	//defer logFile.Close()

	defer func() {
		if err := logFile.Close(); err != nil {
			t.Fatalf("Failed to close log file: %v", err)
		}
	}()

	mdl := builder.LogFunc(func(log string) {
		logFile.WriteString(fmt.Sprintf("%s: %s \n", time.Now().Format(time.RFC3339), log))
	}).Build()

	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/hello", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("Hello log!"))
	})

	server.Start(":8081")
}
