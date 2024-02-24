package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestUserController_GetUser(t *testing.T) {
	g := gin.Default()
	ctrl := &UserController{}
	g.GET("/user/*", ctrl.GetUser)
	g.POST("/user/*", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello %s", "world")
	}) // 这里是 MVC 模式

	g.GET("/static", func(context *gin.Context) {
		// 读文件
		// 谐响应
	})
	_ = g.Run(":8082")

	//http.ListenAndServe(":8083", g) //  Engine 本身可作为一个 Handler 传递到 http 包启动服务

}

// Gin:
// IRoutes:
// Gin 没有 Controller 抽象，但是有 HandlerFunc 抽象 ---> MVC 应该是用户组织 Web 项目的模式，而不是我们中间件设计者要考虑的。
// Engine:
// 实现路由树，提供注册和路由匹配功能
// 本身可作为一个 Handler 传递到 http 包启动服务
// Engine 路由树功能本质上依赖于 methodTree
// HandlerFunc 定义核心抽象 ---> 处理逻辑(业务代码)
// HandlersChain 构造责任链模式
// Context 定义核心抽象 ---> 请求上下文 (提供 API)
