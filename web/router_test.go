package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

// TestRouter_addRoute() 测试注册路由
func TestRouter_addRoute(t *testing.T) {
	// 1.构造路由树
	// 2.验证路由树

	mockHandler := func(ctx Context) {}

	type fields struct {
		trees map[string]*node
	}

	type args struct {
		method     string
		path       string
		handleFunc HandleFunc
	}

	trueTests := []struct {
		name       string
		fields     fields
		args       args
		wantRouter *router
	}{
		// 测试 GET 方法
		{ // 根节点需要特殊处理
			name:       "GET /",
			fields:     fields{trees: make(map[string]*node)},
			args:       args{method: http.MethodGet, path: "/", handleFunc: mockHandler},
			wantRouter: &router{trees: map[string]*node{http.MethodGet: {path: "/", handler: mockHandler}}},
		},
		{
			name:   "GET /user",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/user", handleFunc: mockHandler},
			wantRouter: &router{trees: map[string]*node{http.MethodGet: {path: "/", children: map[string]*node{
				"user": {path: "user", handler: mockHandler},
			}}}},
		},
		{
			name:   "GET /user/home",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/user/home", handleFunc: mockHandler},
			wantRouter: &router{trees: map[string]*node{http.MethodGet: {path: "/", children: map[string]*node{
				"user": {path: "user", children: map[string]*node{"home": {path: "home", handler: mockHandler}}},
			}}}},
		},
		{
			name:   "GET /order/detail",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/order/detail", handleFunc: mockHandler},
			wantRouter: &router{trees: map[string]*node{http.MethodGet: {path: "/", children: map[string]*node{
				"order": {path: "order", children: map[string]*node{"detail": {path: "detail", handler: mockHandler}}},
			}}}},
		},
		// 测试 POST 方法
		{
			name:   "POST /order/create",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodPost, path: "/order/create", handleFunc: mockHandler},
			wantRouter: &router{trees: map[string]*node{http.MethodPost: {path: "/", children: map[string]*node{
				"order": {path: "order", children: map[string]*node{"create": {path: "create", handler: mockHandler}}},
			}}}},
		},
		{
			name:   "POST /login",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodPost, path: "/login", handleFunc: mockHandler},
			wantRouter: &router{trees: map[string]*node{http.MethodPost: {path: "/", children: map[string]*node{
				"login": {path: "login", handler: mockHandler},
			}}}},
		},
		//{ // 不支持前导没有 "/" ---> router 加校验
		//	method: http.MethodPost,
		//	path:   "login",
		//}
	}

	for _, tt := range trueTests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				trees: tt.fields.trees,
			}
			r.addRoute(tt.args.method, tt.args.path, tt.args.handleFunc)

			// 不能直接断言, 因为 HandleFunc 不是可比较的类型
			// assert.Equal(t, wantRouter, r)
			msg, ok := tt.wantRouter.equal(r)
			assert.True(t, ok, msg)
		})
	}

	// 非法用例
	r := newRouter()
	falseTests := []struct {
		name   string
		fields fields
		args   args
		//wantRouter router
		wantErr string
	}{
		{
			name:   "空字符串",
			fields: fields{trees: make(map[string]*node)},
			args: args{
				method:     http.MethodGet,
				path:       "",
				handleFunc: mockHandler,
			},
			wantErr: "web:路径不能为空字符串",
		},
		{
			name:   "前导没有 /",
			fields: fields{trees: make(map[string]*node)},
			args: args{
				method:     http.MethodGet,
				path:       "a/b/c",
				handleFunc: mockHandler,
			},
			wantErr: "web:路径必须以 / 开头",
		},
		{
			name:   "后缀有 /",
			fields: fields{trees: make(map[string]*node)},
			args: args{
				method:     http.MethodGet,
				path:       "/a/b/c/",
				handleFunc: mockHandler,
			},
			wantErr: "web:路径不能以 / 结尾",
		},
		{
			name:   "路由包含连续的 /",
			fields: fields{trees: make(map[string]*node)},
			args: args{
				method:     http.MethodGet,
				path:       "/a//b/c",
				handleFunc: mockHandler,
			},
			wantErr: "web:路径不能包含连续的 / ",
		},
		{
			name: "根节点重复注册",
			fields: fields{
				trees: map[string]*node{http.MethodGet: {path: "/", handler: mockHandler}},
			},
			args: args{
				method:     http.MethodGet,
				path:       "/",
				handleFunc: mockHandler,
			},
			wantErr: "web: 路由冲突, 重复注册 [/] ",
		},
		// TODO: 子节点重复注册
		//{
		//	name: "子节点重复注册",
		//	fields: fields{
		//		trees: map[string]*node{
		//			http.MethodGet: {
		//				path: "/",
		//				children: map[string]*node{
		//					"a": {
		//						path: "a",
		//						children: map[string]*node{
		//							"b": {
		//								path: "b",
		//								children: map[string]*node{
		//									"c": {
		//										path:    "c",
		//										handler: mockHandler,
		//									},
		//								},
		//							},
		//						},
		//					},
		//				},
		//			},
		//		},
		//	},
		//	args: args{
		//		method:     http.MethodGet,
		//		path:       "/a/b/c",
		//		handleFunc: mockHandler,
		//	},
		//	wantErr: "web: 路由冲突, 重复注册 [/a/b/c] ",
		//},
	}

	for _, ft := range falseTests {
		t.Run(ft.name, func(t *testing.T) {
			r.trees = ft.fields.trees
			assert.PanicsWithValue(t, ft.wantErr, func() {
				r.addRoute(ft.args.method, ft.args.path, ft.args.handleFunc)
			})
		})
	}
}

// TestRouter_findRoute() 测试查找路由
func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		// 测试用例
		// 注册路由
		{
			method: http.MethodDelete,
			path:   "/",
		},
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
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
	}{
		{ // 方法不存在
			name:   "method not found",
			method: http.MethodOptions,
			path:   "/order/detail",
			//wantFound: false,
		},
		{ // 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			wantNode: &node{
				handler: mockHandler,
				path:    "detail",
			},
		},
		{ // 命中了, 但是没有 handler
			name:      "order",
			method:    http.MethodGet,
			path:      "/order",
			wantFound: true,
			wantNode: &node{
				//handler: mockHandler,
				path: "order",
				children: map[string]*node{
					"detail": {
						handler: mockHandler,
						path:    "detail",
					},
				},
			},
		},
		{ // 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			msg, ok := tc.wantNode.equal(n)
			assert.True(t, ok, msg)
		})
	}

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
