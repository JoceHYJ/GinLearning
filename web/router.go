package web

import "strings"

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

// AddRoute 添加路由的方法
func (r *router) AddRoute(method, path string, handleFunc HandleFunc) {
	// 找到对应的树
	root, ok := r.trees[method]
	if !ok {
		// 没有根节点则创建
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	// 切割 path
	// /user/home 会被切割成 ["", "user", "home"]三段
	// 第一段是空的，需要去掉第一段
	path = path[1:]
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		// 递归寻找位置 --> children
		// 如果中途有节点不存在则创建
		children := root.childOrCreate(seg)
		root = children
	}
	// 把 handler 挂载到 root 上(赋值)
	root.handler = handleFunc
}

// childOrCreate 用于查找或创建节点的子节点
func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 如果子节点不存在，则新建一个
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

type node struct {
	path string

	//children []*node
	// 子 path 到子节点的映射
	children map[string]*node

	// 缺少代表用户注册的业务逻辑
	handler HandleFunc
}
