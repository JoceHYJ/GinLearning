package beego

import "github.com/beego/beego/v2/server/web"

// beego 依赖: go get github.com/beego/beego/v2@latest

// UserController 控制器: beego 中要求组合 web.Controller
type UserController struct {
	web.Controller
}

func (c *UserController) GetUser() {
	c.Ctx.WriteString("hello tomato")
}

func (c *UserController) CreateUser() {
	u := &User{}
	err := c.Ctx.BindJSON(u)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	_ = c.Ctx.JSONResp(u)
}

type User struct {
	Name string
}
