package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

// TestRouter_addRoute() 测试路由树
func TestRouter_addRoute(t *testing.T) {
	// 1.构造路由树
	// 2.验证路由树
	testRoutes := []struct {
		method  string
		path    string
		handler HandleFunc
	}{
		// 测试用例

		// 测试 GET 方法
		{ // 根节点需要特殊处理
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail", // 没有注册 handler 的节点 --> order
		},

		// 测试 POST 方法
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		//{ // 不支持前导没有 "/" ---> router 加校验
		//	method: http.MethodPost,
		//	path:   "login",
		//},
	}

	var mockHandler HandleFunc = func(ctx Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 3.断言两者相等
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler, // 增加了根节点测试用例，所以添加 handler
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler, // 增加了 /user 的测试用例
						children: map[string]*node{
							"home": &node{
								path:     "home",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						// 不需要 handler
						children: map[string]*node{
							"detail": &node{
								path:     "detail",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:     "create",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
					"login": &node{
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	// 不能直接断言, 因为 HandleFunc 不是可比较的类型
	// assert.Equal(t, wantRouter, r)

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 增加测试用例: path 的格式不符合要求
	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web:路径不能为空字符串")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "web:路径不能以 / 结尾")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a//c", mockHandler)
	}, "web:路径不能包含连续的 / ")

	// 增加测试用例: 路由重复注册
	// 根节点路径重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "web: 路由冲突, 重复注册 [/] ")

	// 子节点路径重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "web: 路由冲突, 重复注册 [/a/b/c] ")

}

// equal 方法: 断言
// string 返回错误信息 --> 排查问题
// bool 返回是否相等
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		// v, dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不匹配"), false
	}

	// 比较 handler --> 利用反射
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("Handler 不匹配"), false
	}

	// 递归比较子节点
	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "匹配成功", true
}
