package web

import (
	"fmt"
	"regexp"
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
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
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

// addRoute 注册路由
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
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
	// 全新的 HTTP 方法，创建根节点
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
		root.route = "/"
		return
	}

	// 切割 path
	// /user/home 会被切割成 ["", "user", "home"]三段
	// 第一段是空的，需要去掉第一段
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web:路径不能包含连续的 / ") // 不能有 //a/b, /a//b 之类的路由
		}
		// 递归寻找位置 --> children
		// 如果中途有节点不存在则创建
		child := root.childOrCreate(seg)
		root = child
		//root = root.childOrCreate(seg)
	}
	// 避免子节点路径重复注册
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突, 重复注册 [%s]", path))
	}
	// 把 handler 挂载到 root 上(赋值)
	root.handler = handleFunc
	root.route = path
}

// findRoute 查找路由的方法(查找对应的节点)
// 返回的 node 内部 Handler 不为 nil --> 路由注册成功
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
	////pathParams := make(map[string]string)
	//var pathParams map[string]string
	//for _, seg := range segs {
	//	child, paramChild, found := root.childOf(seg)
	//	if !found {
	//		return nil, false
	//	}
	//	root = child
	//	// 命中了路径参数
	//	if paramChild {
	//		if pathParams == nil {
	//			pathParams = make(map[string]string)
	//		}
	//		// path 是 :id 的形式
	//		pathParams[child.path[1:]] = seg
	//	}
	//}
	//return &matchInfo{
	//	n:          root,
	//	pathParams: pathParams,
	//}, true // 返回找到的节点 ---> 但是不能返回用户是否注册了 handler
	////return root, root.handler != nil // 返回用户是否注册了 handler
	mi := &matchInfo{}
	for _, s := range segs {
		var child *node
		child, ok = root.childOf(s)
		if !ok {
			if root.typ == nodeTypeAny {
				mi.n = root
				return mi, true
			}
			return nil, false
		}
		if child.paramName != "" {
			mi.addValue(child.paramName, s)
		}
		root = child
	}
	mi.n = root
	return mi, true
}

// childOrCreate 用于查找或创建节点的子节点
// 首先会判断 path 是不是通配符路径
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOrCreate(seg string) *node {
	// 通配符匹配
	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册参数路径和通配符匹配, 已有参数路径匹配")
		}

		if n.regChild != nil {
			panic("web: 不允许同时注册正则路由和通配符路由, 已有正则路由")
		}

		if n.starChild == nil {
			n.starChild = &node{
				path: seg,
				typ:  nodeTypeAny,
			}
		}
		return n.starChild
	}

	//// 参数路径匹配
	//if seg[0] == ':' {
	//	if n.starChild != nil {
	//		panic("web: 不允许同时注册参数路径和通配符匹配, 已有通配符匹配")
	//	}
	//
	//	if n.paramChild != nil {
	//		if n.paramChild.path != seg {
	//			panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, seg))
	//		}
	//	} else {
	//		n.paramChild = &node{
	//			path: seg,
	//		}
	//	}
	//
	//	return n.paramChild
	//}

	// 以 : 开头需要进一步解析，判断是正则匹配还是参数路径匹配
	if seg[0] == ':' {
		paramName, expr, isReg := n.parseParam(seg)
		if isReg {
			return n.childOrCreateReg(seg, expr, paramName)
		}
		return n.childOrCreateParam(seg, paramName)
	}

	if n.children == nil {
		//n.children = map[string]*node{}
		n.children = make(map[string]*node)
	}
	child, ok := n.children[seg]
	if !ok {
		// 如果子节点不存在，则新建一个
		child = &node{
			path: seg,
			typ:  nodeTypeStatic,
		}
		n.children[seg] = child
	}
	return child
}

func (n *node) childOrCreateParam(path string, paramName string) *node {
	if n.regChild != nil {
		panic(fmt.Sprintf("web: 不允许同时注册正则路由和参数路由, 已有正则路由"))
	}
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 不允许同时注册参数路径和通配符匹配, 已有通配符匹配"))
	}
	if n.paramChild != nil {
		if n.paramChild.path != path {
			panic(fmt.Sprintf("web: 路由冲突, 参数路由冲突, 已有 %s, 新注册 %s", n.paramChild.path, path))
		}
	} else {
		n.paramChild = &node{path: path, paramName: paramName, typ: nodeTypeParam}
	}
	return n.paramChild
}

func (n *node) childOrCreateReg(path string, expr string, paramName string) *node {
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 不允许同时注册正则路由和通配符路由, 已有通配符路由"))
	}
	if n.paramChild != nil {
		panic(fmt.Sprintf("web: 不允许同时注册正则路由和参数路由, 已有参数路由"))
	}
	if n.regChild != nil {
		if n.regChild.regExpr.String() != expr || n.paramName != paramName {
			panic(fmt.Sprintf("web: 路由冲突, 正则路由冲突, 已有 %s, 新注册 %s", n.regChild.path, path))
		}
	} else {
		regExpr, err := regexp.Compile(expr)
		if err != nil {
			panic(fmt.Errorf("web: 正则表达式错误 %w", err))
		}
		n.regChild = &node{path: path, paramName: paramName, regExpr: regExpr, typ: nodeTypeReg}
	}
	return n.regChild
}

// childOf 用于查找子节点
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return n.childOfNonStatic(path)
	}
	res, ok := n.children[path]
	if !ok {
		return n.childOfNonStatic(path)
	}
	return res, ok
}

// 优先考虑静态匹配
// 匹配失败则尝试通配符匹配
// 参数路径匹配
// 返回值参数: 第一个是子节点，第二个标记是否是路径参数，第三个标记是否命中
//	func (n *node) childOf(path string) (*node, bool, bool) {
//		if n.children == nil {
//			//return nil, false
//			if n.paramChild != nil {
//				return n.paramChild, true, true
//			}
//			return n.starChild, false, n.starChild != nil
//		}
//		child, ok := n.children[path]
//		if !ok {
//			if n.paramChild != nil {
//				return n.paramChild, true, true
//			}
//			return n.starChild, false, n.starChild != nil
//		}
//		return child, false, ok
//	}

// childOfNonStatic 从非静态匹配的子节点里面查找
func (n *node) childOfNonStatic(path string) (*node, bool) {
	if n.regChild != nil {
		if n.regChild.regExpr.Match([]byte(path)) {
			return n.regChild, true
		}
	}
	if n.paramChild != nil {
		return n.paramChild, true
	}
	return n.starChild, n.starChild != nil
}

// parseParam 用于解析判断是不是正则表达式
// 第一个返回值是参数名字
// 第二个返回值是正则表达式
// 第三个返回值为 true 则说明是正则路由
func (n *node) parseParam(seg string) (string, string, bool) {
	// 去除 :
	seg = seg[1:]
	segs := strings.SplitN(seg, "(", 2)
	if len(segs) == 2 {
		expr := segs[1]
		if strings.HasSuffix(expr, ")") {
			return segs[0], expr[:len(expr)-1], true
		}
	}
	return seg, "", false
}

type nodeType int

const (
	// 静态匹配
	nodeTypeStatic nodeType = iota
	// 正则匹配
	nodeTypeReg
	// 参数匹配
	nodeTypeParam
	// 通配符匹配
	nodeTypeAny
)

type node struct {
	typ nodeType

	path string

	// 静态匹配的节点
	//children []*node
	// 子 path 到子节点的映射
	children map[string]*node

	// 通配符匹配节点
	starChild *node

	// 参数路径匹配节点
	paramChild *node

	// 正则路由和参数路由都会使用该字段
	paramName string

	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp

	// handler 命中路由后执行的逻辑
	handler HandleFunc

	// route 到达该节点的完整路由路径
	route string
}

// 参数匹配信息
type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
