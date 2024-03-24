package session

import (
	"context"
	"net/http"
)

// Session 接口
type Session interface {
	// Get 获取 session 值
	Get(ctx context.Context, key string) (string, error)
	// Set 设置 session 值
	Set(ctx context.Context, key string, val string) error
	// ID 获取 session id
	ID() string
}

// Store 管理 Session
// 从设计的角度来说，Generate 方法和 Refresh 在处理 Session 过期时间上有点关系
// 也就是说，如果 Generate 设计为接收一个 expiration 参数，
// 那么 Refresh 也应该接收一个 expiration 参数。
// 因为这意味着用户来管理过期时间
type Store interface {
	// Generate 生成新的 session
	Generate(ctx context.Context, id string) (Session, error)
	// Refresh 刷新 session 的过期时间
	// 这种设计是一直用同一个 id 的
	// 如果想支持 Refresh 换 ID，那么可以重新生成一个，并移除原有的
	// 又或者 Refresh(ctx context.Context, id string) (Session, error)
	// 其中返回的是一个新的 Session
	Refresh(ctx context.Context, id string) error
	// Remove 移除 session
	Remove(ctx context.Context, id string) error
	//Get 获取 session
	Get(ctx context.Context, id string) (Session, error)
}

type Propagator interface {
	// Inject 将 session id 注入到 context 中
	// Inject 必须是幂等的，也就是说多次调用应该和一次调用效果一样
	Inject(id string, writer http.ResponseWriter) error
	// Extract 从 context 中提取 session id
	// e.g. cookie 将从 session 中提取出来
	Extract(req *http.Request) (string, error)
	// Remove 从 http.ResponseWriter 中移除 session id
	Remove(writer http.ResponseWriter) error
}
