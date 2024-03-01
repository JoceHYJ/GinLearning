package web

// 全静态匹配

// router
// 用来支持对路由树的操作
// 代表路由树(森林)
type router struct {
	// Beego Gin: HTTP method 对应一棵树
	// GET POST 也各对应一棵树
	//trees map[string]tree

	// http method => 路由树根节点
	trees map[string]*node
}

//type tree struct {
//	root *node
//}

// newRouter 创建路由的方法
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRoute(method, path string, handler HandleFunc) {

}

type node struct {
	path string

	//children []*node
	// 子 path 到子节点的映射
	children map[string]*node

	// 缺少代表用户注册的业务逻辑
	handler HandleFunc
}
