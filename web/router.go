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
func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// AddRoute 添加路由的方法
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
			panic("web: 路由冲突, 重复注册 [/] ")
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
		children := root.childOrCreate(seg)
		root = children
	}
	// 避免子节点路径重复注册
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突, 重复注册 [%s]", path))
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
