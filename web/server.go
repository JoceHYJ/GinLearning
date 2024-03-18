package web

import (
	"log"
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

	mdls      []Middleware
	tplEngine TemplateEngine
}

// Option 模式

type HTTPServerOption func(server *HTTPServer)

// NewHTTPServer 初始化,创建一个 HTTPServer (路由器)实例
func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	//return &HTTPServer{
	//	router: newRouter(),
	//}
	res := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

func ServerWithTemplateEngine(tplEngine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tplEngine
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
		Resp:      writer,
		Req:       request,
		tplEngine: h.tplEngine,
	}
	// 最后一个应该是 HTTPServer 执行路由匹配，执行用户代码
	root := h.serve
	// 从后往前回溯组装链条
	// 把后一个作为前一个的 next 构造链条
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}
	//h.serve(ctx)
	// 从前往后
	// 第一个应该是回写响应的
	// 因为在调用 next 之后才会回写响应
	// 所以实际上 flashResp 应该是最后一个步骤
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}
	root = m(root)
	root(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	// 查找路由, 并执行命中的业务逻辑
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	//if !ok || info.n.handler == nil {
	//	ctx.Resp.WriteHeader(http.StatusNotFound)
	//	_, _ = ctx.Resp.Write([]byte("404 NOT FOUND"))
	//	return
	//}

	if !ok || info.n == nil || info.n.handler == nil {
		ctx.RespStatusCode = http.StatusNotFound
		return
	}
	ctx.PathParams = info.pathParams
	//ctx.MatchedRoute = info.n.path
	ctx.MatchedRoute = info.n.route
	info.n.handler(ctx)
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

// Start 启动服务器, 用户定义指定端口
// 编程接口
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

// 配置接口
// 问题:
// 路径: 相对 or 绝对
// 配置文件格式: json yaml xml
//func NewHTTPServerV1(cfgFilePath string) *HTTPServer {
//	// 加载配置文件
//	// 解析
//    // 初始化
//    return &HTTPServer{}
//}

// Use 向 HTTPServer 实例中添加中间件
func (h *HTTPServer) Use(mdls ...Middleware) {
	if h.mdls == nil {
		h.mdls = mdls
		return
	}
	h.mdls = append(h.mdls, mdls...)
}

// flashResp 回写响应
func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败:", err)
	}
}
