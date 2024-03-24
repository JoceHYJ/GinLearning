package session

import (
	"GinLearning/web"
)

type Manager struct {
	Store
	Propagator
	SessCtxKey string
}

// GetSession 从 ctx 中获取 session
// 成功返回 session 实例缓存到 ctx 的 UserValues 字段中
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}

	val, ok := ctx.UserValues[m.SessCtxKey]
	if ok {
		return val.(Session), nil
	}
	id, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	session, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.SessCtxKey] = session
	return session, nil
}

// InitSession 初始化 Session 并注入到 http 的 response 中
func (m *Manager) InitSession(ctx *web.Context, id string) (Session, error) {
	session, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return session, nil
}

// RefreshSession 刷新 session 的过期时间
func (m *Manager) RefreshSession(ctx *web.Context) (Session, error) {
	session, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	// 刷新 session 过期时间
	if err = m.Refresh(ctx.Req.Context(), session.ID()); err != nil {
		return nil, err
	}
	// 重新注入 session
	if err = m.Inject(session.ID(), ctx.Resp); err != nil {
		return nil, err
	}
	return session, nil
}

// RemoveSession 从 ctx 中移除 session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	if err = m.Store.Remove(ctx.Req.Context(), session.ID()); err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Resp)
}
