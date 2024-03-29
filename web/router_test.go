package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

// 该单元测试没有支持可路由的中间件测试, 测试代码通过需要将 web 框架退回上一版本

// TestRouter_addRoute() 测试注册路由
func TestRouter_addRoute(t *testing.T) {
	// 1.构造路由树
	// 2.验证路由树

	mockHandler := func(ctx *Context) {}

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
		wantRouter router
	}{
		// 1.全静态匹配
		{ // 根节点需要特殊处理
			name:   "GET /",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", handler: mockHandler, typ: nodeTypeStatic,
			}}},
		},
		{
			name:   "GET /user",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/user", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"user": {path: "user", handler: mockHandler, typ: nodeTypeStatic},
				}}}},
		},
		{
			name:   "GET /user/home",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/user/home", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"user": {path: "user", typ: nodeTypeStatic, children: map[string]*node{
						"home": {path: "home", handler: mockHandler, typ: nodeTypeStatic}}},
				}}}},
		},
		{
			name:   "GET /order",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/order", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"order": {path: "order", handler: mockHandler, typ: nodeTypeStatic},
				}}}},
		},
		{
			name:   "GET /order/detail",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/order/detail", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"order": {path: "order", typ: nodeTypeStatic, children: map[string]*node{
						"detail": {path: "detail", handler: mockHandler, typ: nodeTypeStatic}}},
				}}}},
		},
		// 测试 POST 方法
		{
			name:   "POST /order/create",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodPost, path: "/order/create", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodPost: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"order": {path: "order", typ: nodeTypeStatic, children: map[string]*node{
						"create": {path: "create", handler: mockHandler, typ: nodeTypeStatic}}},
				}}}},
		},
		{
			name:   "POST /login",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodPost, path: "/login", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodPost: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"login": {path: "login", handler: mockHandler, typ: nodeTypeStatic},
				}}}},
		},
		//{ // 不支持前导没有 "/" ---> router 加校验
		//	method: http.MethodPost,
		//	path:   "login",
		//}

		// 2.通配符匹配
		{
			name:   "GET /order/*",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/order/*", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"order": {path: "order", typ: nodeTypeStatic, starChild: &node{
						path: "*", handler: mockHandler, typ: nodeTypeAny}},
				}}}},
		},
		{
			name:   "GET /*",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/*", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, starChild: &node{
					path: "*", handler: mockHandler, typ: nodeTypeAny,
				}}}},
		},
		{
			name:   "GET /*/*",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/*/*", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, starChild: &node{
					path: "*", typ: nodeTypeAny, starChild: &node{
						path: "*", handler: mockHandler, typ: nodeTypeAny,
					}}}}},
		},
		{
			name:   "GET /*/abc",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/*/abc", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, starChild: &node{
					path: "*", typ: nodeTypeAny, children: map[string]*node{
						"abc": {path: "abc", handler: mockHandler, typ: nodeTypeStatic},
					}}}}},
		},
		{
			name:   "GET /*/abc/*",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/*/abc/*", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, starChild: &node{
					path: "*", typ: nodeTypeAny, children: map[string]*node{
						"abc": {path: "abc", typ: nodeTypeStatic, starChild: &node{
							path: "*", handler: mockHandler, typ: nodeTypeAny}},
					}}}}},
		},

		// 3. 参数路径匹配 eg: /user/:id -> /user/123, id = 123
		{
			name:   "GET /order/detail/:id",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/order/detail/:id", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"order": {path: "order", typ: nodeTypeStatic, children: map[string]*node{"detail": {
						path: "detail", typ: nodeTypeStatic, paramChild: &node{
							path: ":id", handler: mockHandler, typ: nodeTypeParam,
						}}}},
				}}}},
		},
		{
			name:   "GET /param/:id",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/param/:id", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"param": {path: "param", typ: nodeTypeStatic, paramChild: &node{
						path: ":id", handler: mockHandler, typ: nodeTypeParam}},
				}}}},
		},
		{
			name:   "GET /param/:id/detail",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/param/:id/detail", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"param": {path: "param", typ: nodeTypeStatic, paramChild: &node{
						path: ":id", typ: nodeTypeParam, children: map[string]*node{
							"detail": {path: "detail", handler: mockHandler, typ: nodeTypeStatic},
						}}},
				}}}},
		},
		//{ // 已有静态匹配增加通配符匹配
		//	name: "GET /param/:id/(detail)*",
		//	fields: fields{trees: map[string]*node{http.MethodGet: {
		//		path: "/", children: map[string]*node{
		//			"param": {path: "param", paramChild: &node{
		//				path: ":id", children: map[string]*node{
		//					"detail": {path: "detail", handler: mockHandler},
		//				}}},
		//		}}}},
		//	args: args{method: http.MethodGet, path: "/param/:id/*", handleFunc: mockHandler},
		//	wantRouter: router{trees: map[string]*node{http.MethodGet: {path: "/", children: map[string]*node{
		//		"param": {path: "param", paramChild: &node{
		//			path: ":id", starChild: &node{
		//				path: "*", handler: mockHandler}, children: map[string]*node{
		//				"detail": {path: "detail", handler: mockHandler},
		//			},
		//		},
		//		}}},
		//	}},
		//},
		{ // param, star 同时
			name:   "GET /param/:id/*",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodGet, path: "/param/:id/*", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"param": {path: "param", typ: nodeTypeStatic, paramChild: &node{
						path: ":id", typ: nodeTypeParam, starChild: &node{
							path: "*", handler: mockHandler, typ: nodeTypeAny}},
					}}},
			}},
		},
		// 正则匹配
		{
			name:   "DELETE /reg/:id(.*)",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodDelete, path: "/reg/:id(.*)", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodDelete: {
				path: "/", typ: nodeTypeStatic, children: map[string]*node{
					"reg": {path: "reg", typ: nodeTypeStatic, regChild: &node{
						path: ":id(.*)", paramName: ":id", handler: mockHandler, typ: nodeTypeReg,
					}}},
			}}},
		},
		{
			name:   "DELETE /:name(^.+$)/abc",
			fields: fields{trees: make(map[string]*node)},
			args:   args{method: http.MethodDelete, path: "/:name(^.+$)/abc", handleFunc: mockHandler},
			wantRouter: router{trees: map[string]*node{http.MethodDelete: {
				path: "/", typ: nodeTypeStatic, regChild: &node{
					path: ":name(^.+$)", paramName: "name", typ: nodeTypeReg, children: map[string]*node{
						"abc": {path: "abc", handler: mockHandler, typ: nodeTypeStatic},
					},
				},
			}}},
		},
	}

	for _, tt := range trueTests {
		t.Run(tt.name, func(t *testing.T) {
			r := router{
				trees: tt.fields.trees,
			}
			r.addRoute(tt.args.method, tt.args.path, tt.args.handleFunc)

			// 不能直接断言, 因为 HandleFunc 不是可比较的类型
			// assert.Equal(t, wantRouter, r)
			msg, ok := tt.wantRouter.equal(&r)
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
		// 1.全静态匹配
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
				trees: map[string]*node{http.MethodGet: {
					path: "/", handler: mockHandler,
				}},
			},
			args: args{
				method:     http.MethodGet,
				path:       "/",
				handleFunc: mockHandler,
			},
			wantErr: "web: 路由冲突, 重复注册 [/]",
		},
		{
			name: "子节点重复注册",
			fields: fields{
				trees: map[string]*node{http.MethodGet: {
					path: "/", children: map[string]*node{
						"a": {path: "a", children: map[string]*node{
							"b": {path: "b", children: map[string]*node{
								"c": {path: "c", handler: mockHandler},
							}}}},
					},
				}},
			},
			args: args{
				method:     http.MethodGet,
				path:       "/a/b/c",
				handleFunc: mockHandler,
			},
			wantErr: "web: 路由冲突, 重复注册 [/a/b/c]",
		},
		{
			name: "不允许同时注册参数路径和通配符匹配,已有通配符匹配",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", starChild: &node{
						path: "*", handler: mockHandler,
					}},
				},
			}}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/:id",
				handleFunc: mockHandler,
			},
			wantErr: "web: 不允许同时注册参数路径和通配符匹配, 已有通配符匹配",
		},
		{
			name: "不允许同时注册参数路径和通配符匹配,已有参数路径匹配",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", paramChild: &node{
						path: ":id", handler: mockHandler,
					}},
				},
			}}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/*",
				handleFunc: mockHandler,
			},
			wantErr: "web: 不允许同时注册参数路径和通配符匹配, 已有参数路径匹配",
		},
		{
			name: "不允许同时注册正则路由和通配符路由, 已有通配符路由",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", starChild: &node{
						path: "*", handler: mockHandler,
					}},
				},
			}}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/:id(.*)",
				handleFunc: mockHandler,
			},
			wantErr: "web: 不允许同时注册正则路由和通配符路由, 已有通配符路由",
		},
		{
			name: "不允许同时注册正则路由和参数路由, 已有参数路由",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", children: map[string]*node{
						"b": {path: "b", paramChild: &node{
							path: ":id", handler: mockHandler,
						}},
					}},
				},
			}}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/b/:id(.*)",
				handleFunc: mockHandler,
			},
			wantErr: "web: 不允许同时注册正则路由和参数路由, 已有参数路由",
		},
		{
			name: "不允许同时注册正则路由和通配符路由, 已有正则路由",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", children: map[string]*node{
						"b": {path: "b", regChild: &node{
							path: ":id(.*)", handler: mockHandler,
						},
						},
					}},
				}},
			}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/b/*",
				handleFunc: mockHandler,
			},
			wantErr: "web: 不允许同时注册正则路由和通配符路由, 已有正则路由",
		},
		{
			name: "不允许同时注册正则路由和参数路由, 已有正则路由",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", children: map[string]*node{
						"b": {path: "b", regChild: &node{
							path: ":id(.*)", handler: mockHandler,
						},
						},
					}},
				}},
			}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/b/:id",
				handleFunc: mockHandler,
			},
			wantErr: "web: 不允许同时注册正则路由和参数路由, 已有正则路由",
		},
		{
			name: "路由冲突, 参数路由冲突",
			fields: fields{trees: map[string]*node{http.MethodGet: {
				path: "/", children: map[string]*node{
					"a": {path: "a", children: map[string]*node{
						"b": {path: "b", children: map[string]*node{
							"c": {path: "c", paramChild: &node{
								path: ":id", handler: mockHandler,
							}},
						}},
					}},
				},
			}}},
			args: args{
				method:     http.MethodGet,
				path:       "/a/b/c/:name",
				handleFunc: mockHandler,
			},
			wantErr: "web: 路由冲突, 参数路由冲突, 已有 :id, 新注册 :name",
		},
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

	type fields struct {
		trees map[string]*node
	}

	type args struct {
		method string
		path   string
	}

	testRoute := []struct {
		method string
		path   string
	}{
		// 测试用例
		// 1.全静态匹配
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
		//{
		//	method: http.MethodGet,
		//	path:   "/user",
		//},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		//{
		//	method: http.MethodGet,
		//	path:   "/user/home",
		//},
		{
			method: http.MethodGet,
			path:   "/order/detail", // 没有注册 handler 的节点 --> order
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		{
			method: http.MethodPost,
			path:   "/login/:username",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:id([0-9]+)/home",
		},
	}

	mockHandler := func(ctx *Context) {}

	r := newRouter()
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	test := []struct {
		name   string
		fields fields
		args   args
		//wantNode  *node
		wantMatchInfo *matchInfo
		wantFound     bool
	}{
		// 1.全静态匹配
		{ // 方法不存在
			name:      "method not found",
			args:      args{method: http.MethodOptions, path: "/order/detail"},
			wantFound: false,
		},
		{
			name:      "path not found",
			args:      args{method: http.MethodGet, path: "/abc"},
			wantFound: false,
		},
		{ // 完全命中
			name:      "order detail",
			args:      args{method: http.MethodGet, path: "/order/detail"},
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				handler: mockHandler,
				path:    "detail",
			}},
		},
		{ // 命中了, 但是没有 handler
			name:      "order",
			args:      args{method: http.MethodGet, path: "/order"},
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				path: "order",
				children: map[string]*node{
					"detail": {
						handler: mockHandler,
						path:    "detail",
					}}},
			},
		},
		{ // 根节点
			name:      "root",
			args:      args{method: http.MethodDelete, path: "/"},
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				handler: mockHandler,
				path:    "/",
			}},
		},
		// 通配符匹配
		{ // /order/*
			name:      "star match",
			args:      args{method: http.MethodGet, path: "/order/del"},
			wantFound: true,
			wantMatchInfo: &matchInfo{n: &node{
				handler: mockHandler,
				path:    "*",
			}},
		},
		{ // /user/*/home
			name:      "star in middle",
			args:      args{method: http.MethodGet, path: "/user/tomato/home"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "home",
				},
			},
		},
		{
			name: "overflow",
			args: args{method: http.MethodPost, path: "/order/del/sprite"},
		},
		{ // /login/:username
			name:      "login username",
			args:      args{method: http.MethodPost, path: "/login/tomato"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    ":username",
				},
				pathParams: map[string]string{
					"username": "tomato",
				},
			},
		},
		{ // /param/:id
			name:      ":id",
			args:      args{method: http.MethodGet, path: "/param/123"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    ":id",
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{ // /param/:id/*
			name:      ":id/*",
			args:      args{method: http.MethodGet, path: "/param/234/abc"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "*",
				},
				pathParams: map[string]string{
					"id": "234",
				},
			},
		},
		{ // /param/:id/detail
			name:      ":id/detail",
			args:      args{method: http.MethodGet, path: "/param/abc/detail"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "detail",
				},
				pathParams: map[string]string{
					"id": "abc",
				},
			},
		},
		{ // /reg/:id(.*)
			name:      ":id(.*)",
			args:      args{method: http.MethodDelete, path: "/reg/123"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    ":id(.*)",
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{ // /:id([0-9]+)/home
			name:      ":id([0-9]+)",
			args:      args{method: http.MethodDelete, path: "/123/home"},
			wantFound: true,
			wantMatchInfo: &matchInfo{
				n: &node{
					path:    "home",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"id": "123",
				},
			},
		},
		{ // 未命中 /:id([0-9]+)/home
			name:      http.MethodDelete,
			args:      args{method: http.MethodDelete, path: "/abc/home"},
			wantFound: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			info, found := r.findRoute(tt.args.method, tt.args.path)
			assert.Equal(t, tt.wantFound, found)
			if !found {
				return
			}
			//msg, ok := tt.wantNode.equal(n)
			assert.Equal(t, tt.wantMatchInfo.pathParams, info.pathParams)
			//msg, ok := tt.wantMatchInfo.n.equal(info.n)
			//assert.True(t, ok, msg)
			n := tt.wantMatchInfo.n
			wantVal := reflect.ValueOf(info.n.handler)
			gotVal := reflect.ValueOf(n.handler)
			assert.Equal(t, wantVal, gotVal)
		})
	}
}

// equal 方法断言
// string 返回错误信息 --> 排查问题
// bool 返回是否相等
// r 为 wantRouter, y 为 gotRouter
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method %s", k), false
		}
		// v, dst 要相等
		msg, equal := v.equal(dst)
		if !equal {
			return k + "-" + msg, equal
		}
	}
	return "", true
}

// equal: HandleFunc 无法直接比较通过 equal 方法比较
// n: wantNode y: gotNode
func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if n.typ != y.typ {
		return fmt.Sprintf("%s 节点类型不相等 x %d, y %d", n.path, n.typ, y.typ), false
	}

	if n.paramName != y.paramName {
		return fmt.Sprintf("%s 节点参数名字不相等 x %s, y %s", n.path, n.paramName, y.paramName), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}

	if len(n.children) == 0 {
		return "", true
	}

	if n.starChild != nil {
		str, ok := n.starChild.equal(y.starChild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.path, str), false
		}
	}
	if n.paramChild != nil {
		str, ok := n.paramChild.equal(y.paramChild)
		if !ok {
			return fmt.Sprintf("%s 路径参数节点不匹配 %s", n.path, str), false
		}
	}

	if n.regChild != nil {
		str, ok := n.regChild.equal(y.regChild)
		if !ok {
			return fmt.Sprintf("%s 路径参数节点不匹配 %s", n.path, str), false
		}
	}

	//if len(n.children) == 0 {
	//	return "", true
	//}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}
