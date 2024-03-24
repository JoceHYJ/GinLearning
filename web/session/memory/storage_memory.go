package memory

import (
	"GinLearning/web/session"
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Store struct {
	// mutex 确保并发安全 同一个 id 不会被多个 goroutine 同时操作
	mutex sync.RWMutex
	// 利用内存缓存管理 session 过期时间
	sessionCache *cache.Cache
	expiration   time.Duration
}

// NewStore 创建一个新的内存存储(Store 实例)
// 可以通过 Option 设计模式，允许用户控制 session 过期检查间隔
func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessionCache: cache.New(expiration, time.Second),
		expiration:   expiration,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}
	s.sessionCache.Set(sess.ID(), sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, ok := s.sessionCache.Get(id)
	if !ok {
		return errors.New("session: session not found")
	}
	s.sessionCache.Set(sess.(*memorySession).ID(), sess, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessionCache.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.Unlock()
	sess, ok := s.sessionCache.Get(id)
	if !ok {
		return nil, errors.New("session: session not found")
	}
	return sess.(*memorySession), nil
}

type memorySession struct {
	mutex sync.RWMutex
	id    string
	data  map[string]string
}

func (s *memorySession) Get(ctx context.Context, key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.data[key]
	if !ok {
		return "", errors.New("session: key not found")
	}
	return val, nil
}

func (s *memorySession) Set(ctx context.Context, key string, val string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = val
	return nil
}

// ID 返回会话的唯一标识符 -> 后续是只读的不需要锁
func (s *memorySession) ID() string {
	return s.id
}
