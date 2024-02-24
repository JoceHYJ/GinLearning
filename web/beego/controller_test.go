package beego

import (
	"github.com/beego/beego/v2/server/web"
	"testing"
)

func TestUserController(t *testing.T) {
	//go func() {
	//	s := web.NewHttpSever()
	//	s.Run(":8082")
	//	// 启动第二个服务 (Server) ，监听 8082 端口 (跟 8081 是隔离的)
	//	// 这里是为了测试，实际中不会启动第二个服务
	//}()
	// 测试UserController的代码
	web.BConfig.CopyRequestBody = true
	c := &UserController{}
	// 虽然要求组合 Controller
	// 但是注册路由需要使用 web 的包方法 Router 而不是 Controller 的方法 Router (Run同理)
	// ControllerRegister 解决路由注册、匹配
	web.Router("/user", c, "get:GetUser")
	// 监听 8081 端口
	web.Run(":8081") // httpServer 作为服务器抽象, 用于管理应用生命周期和资源隔离单位
}

// beego 核心抽象:
// ControllerRegister 基础
// httpServer 服务器抽象
// Context Controller 辅助(提供 API)

// beego 自带 MVC 模式
