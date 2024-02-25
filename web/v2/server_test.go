package v2

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	// Server 定义一: 只组合 http.Handler
	// 用户在使用时只需调用 http.ListenAndServe 即可
	// 与 https 协议完全无缝衔接，但是难以控制生命周期，缺乏控制力，无法优雅退出
	//var h Server
	//var h Server = &HTTPServer{} ---> 接口没有 Get 方法
	h := &HTTPServer{}

	h.AddRoute(http.MethodGet, "/user", func(ctx Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})

	handler1 := func(ctx Context) {
		fmt.Println("处理第一件事")
	}
	handler2 := func(ctx Context) {
		fmt.Println("处理第二件事")
	}

	// 用户可以自己管理多个 handleFunc, 没必要提供多个
	h.AddRoute(http.MethodGet, "/user", func(ctx Context) {
		handler1(ctx)
		handler2(ctx)
	})

	h.Get("/user", func(ctx Context) {

	})

	//h.AddRoute1(http.MethodGet, "/user", handler1, handler2)

	//h.AddRoute1(http.MethodGet, "/user", func(ctx Context) {
	//	fmt.Println("处理第一件事")
	//}, func(ctx Context) {
	//	fmt.Println("处理第二件事")
	//}) // 2 个 handleFunc

	// 用法一: 完全委托给 http 包
	http.ListenAndServe(":8081", h)
	http.ListenAndServeTLS(":443", "", "", h)

	// Server 定义二:  组合 http.Handler 并增加 Start 方法
	// Server 既可以作为 http.Handler 使用，又可以作为独立的实体，管理生命周期
	// 但是，如果用户不希望使用 ListenAndServe 方法，那么 Server 需要提供 HTTPS 的支持
	// 用法二: 手动管理
	h.Start(":8081")

	// 版本一、二都直接耦合了 Go 自带的 http 包，与第三方 http 库 (fasthttp) 进行对接会很困难
}
