package web

import (
	"fmt"
	"strings"
)

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
func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 添加路由的方法
func (r *router) addRoute(method, path string, handleFunc HandleFunc) {
	// 对 path 加限制 --> 只支持 /user/home 这种格式
	if path == "" {
		panic("web:路径不能为空字符串")
	}
	if path[0] != '/' {
		panic("web:路径必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web:路径不能以 / 结尾")
	}
	// 中间包含连续的 // --> 可以 strings.contains("//")
	// 在 seg 时进行处理

	// 找到对应的树
	root, ok := r.trees[method]
	if !ok {
		// 没有根节点则创建
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 根节点需要特殊处理 path
	if path == "/" {
		// 避免根节点路由重复注册
		if root.handler != nil {
			panic("web: 路由冲突, 重复注册 [/]")
		}
		root.handler = handleFunc
		return
	}

	// 切割 path
	// /user/home 会被切割成 ["", "user", "home"]三段
	// 第一段是空的，需要去掉第一段
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web:路径不能包含连续的 / ")
		}
		// 递归寻找位置 --> children
		// 如果中途有节点不存在则创建
		child := root.childOrCreate(seg)
		root = child
	}
	// 避免子节点路径重复注册
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突, 重复注册 [%s]", path))
	}
	// 把 handler 挂载到 root 上(赋值)
	root.handler = handleFunc
}

// findRoute 查找路由的方法
func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	// 沿着树进行 DFS
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}
	// 把前置和后置的 / 都去掉
	path = strings.Trim(path, "/")
	// 按照 / 切割 path
	segs := strings.Split(path, "/")
	//pathParams := make(map[string]string)
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			// path 是 :id 的形式
			pathParams[child.path[1:]] = seg
		}
	}
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true // 返回找到的节点 ---> 但是不能返回用户是否注册了 handler
	//return root, root.handler != nil // 返回用户是否注册了 handler
}

// childOrCreate 用于查找或创建节点的子节点
func (n *node) childOrCreate(seg string) *node {

	// 参数路径匹配
	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: 不允许同时注册参数路径和通配符匹配, 已有通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	// 通配符匹配
	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册参数路径和通配符匹配, 已有参数路径匹配")
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}

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

// childOf 用于查找子节点
// 优先考虑静态匹配
// 匹配失败则尝试通配符匹配
// 参数路径匹配
// 返回值参数: 第一个是子节点，第二个标记是否是路径参数，第三个标记是否命中
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		//return nil, false
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}

type node struct {
	path string

	// 静态匹配的节点
	//children []*node
	// 子 path 到子节点的映射
	children map[string]*node

	// 通配符匹配节点
	starChild *node

	// 参数路径匹配节点
	paramChild *node

	// 缺少代表用户注册的业务逻辑
	handler HandleFunc
}

// 参数匹配信息
type matchInfo struct {
	n          *node
	pathParams map[string]string
}
