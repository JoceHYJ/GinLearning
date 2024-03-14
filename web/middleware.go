package web

// AOP 设计方案

// Middleware 函数式的责任链模式
// 洋葱模式
type Middleware func(next HandleFunc) HandleFunc

// java 中会使用后面两种

// Interceptor
//type MiddlewareV1 interface {
//	Invoke(next HandleFunc) HandleFunc
//}
//type Interceptor interface {
//	Before(ctx *Context)
//	After(ctx *Context)
//	Surround(ctx *Context)
//}

// 类似于 Gin
//type Chain []HandleFuncV1
//
//type HandleFuncV1 func(*Context) (next bool)
//
//type ChainV1 struct {
//	handlers []HandleFuncV1
//}
//
//func (c ChainV1) Run(ctx *Context) {
//	for _, h := range c.handlers {
//		next := h(ctx)
//		// 中断执行
//		if !next {
//			return
//		}
//	}
//}
