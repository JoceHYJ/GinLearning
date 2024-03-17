package errhdl

import "GinLearning/web"

type MiddlewareBuilder struct {
	// 这种设计只能返回固定的值, 不能做到动态渲染
	// 对于不同情况是否需要页面的重定向
	// 需要在实现可路由的 Middleware 之后再进行设计
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: make(map[int][]byte, 64),
	}
}

// RegisterError 将注册一个错误码，并且返回特定的错误数据
// 这个错误数据可以是一个字符串，也可以是一个页面
func (m *MiddlewareBuilder) RegisterError(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m
}

func (m *MiddlewareBuilder) AddCode(code int, resp []byte) *MiddlewareBuilder {
	m.resp[code] = resp
	return m
}

func (m *MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				// 篡改结果
				ctx.RespData = resp
			}
		}
	}
}
