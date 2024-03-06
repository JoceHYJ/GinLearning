package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 确保一定实现了 Server 接口
var _ Server = &HTTPServer{}

// Server 接口定义
type Server interface {
	http.Handler             // 1. 组合 http.Handler
	Start(addr string) error // 2. 组合 http.Handler 并增加 Start 方法
	//Start1() error           // 不接收 addr 参数

	// addRoute 增加路由注册功能
	// method: HTTP 方法
	// path: 请求路径(路由)
	// handleFunc: 处理函数(业务逻辑)
	addRoute(method string, path string, handleFunc HandleFunc)
	// addRoute1 提供多个 handleFunc: 用户自己组合
	//addRoute1(method string, path string, handles ...HandleFunc)
}

type HTTPServer struct {
	// addr string 创建的时候传递, 而不是 Start 接受，都是可以的
	router
	//*router
	// r *router
	// 三种组合方式都是可以的
}

// NewHTTPServer 初始化,创建一个 HTTPServer (路由器)实例
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// 定义二: 如果用户不希望使用 ListenAndServe 方法，那么 Server 需要提供 HTTPS 的支持
//type HTTPSServer struct {
//	HTTPServer
//}

// Web框架核心入口 ServeHTTP
// ServeHTTP -> HTTPServer 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// http.MethodPut
	// Web 框架代码
	// 1.Context 构建
	// 2.路由匹配
	// 3.执行业务逻辑
	ctx := &Context{
		Resp: writer,
		Req:  request,
	}
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	// 查找路由, 并执行命中的业务逻辑
	n, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || n.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("404 NOT FOUND"))
		return
	}
	n.n.handler(ctx)
}

// addRoute 方法
// 只接收一个 HandleFunc: 因为只希望它注册业务逻辑
// addRoute 方法最终会和路由树交互
// 核心 API
//func (h *HTTPServer) addRoute(method string, path string, handleFunc HandleFunc) {
//	// 这里注册到路由树中
//	panic("implement me")
//}

// 衍生 API ---> 都可以委托给 核心 API addRoute 实现

// Get 方法
// 只定义在实现里(HTTPServer)而不定义在接口里 --> 接口小而美
// Get 等核心 API (HTTP 方法注册的) 都委托给 addRoute(Handle) 方法实现
func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

// Post 方法
func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

// Delete 方法
func (h *HTTPServer) Delete(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodDelete, path, handleFunc)
}

// Put 方法
func (h *HTTPServer) Put(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}

// Options 方法
func (h *HTTPServer) Options(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

//....

//addRoute1 方法
//为了通过编译添加
//func (h *HTTPServer) addRoute1(method string, path string, handles ...HandleFunc) {
//}

func (h *HTTPServer) Start(addr string) error {
	// 也可以自己内部创建 Server 来启动服务
	//http.Server{}
	// 用法二: 自己管理生命周期(Listen->Serve)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 区别:(生命周期的回调)
	// after start 回调
	// 往 admin 注册实例
	// 执行业务所需前置条件
	return http.Serve(l, h)
	//panic("implement me")
}
