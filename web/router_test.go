package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

// TestRouter_AddRoute 测试路由树
func TestRouter_AddRoute(t *testing.T) {
	// 1.构造路由树
	// 2.验证路由树
	testRoutes := []struct {
		method  string
		path    string
		handler HandleFunc
	}{
		{ // 测试 GET 方法
			method: http.MethodGet,
			path:   "/user/home",
		},
	}

	var mockHandler HandleFunc = func(ctx Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	// 3.断言两者相等
	//wantRouter := &router{
	//	trees: map[string]*node{
	//		http.MethodGet: &node{
	//			path: "/",
	//			children: map[string]*node{
	//				"user": &node{
	//					path: "user",
	//					children: map[string]*node{
	//						"home": &node{
	//							path:     "home",
	//							children: map[string]*node{},
	//							handler:  mockHandler,
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {
								path:     "home",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
				},
			},
		},
	}

	// 不能直接断言, 因为 HandleFunc 不是可比较的类型
	// assert.Equal(t, wantRouter, r)

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)
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
