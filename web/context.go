package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Context struct {
	Resp       http.ResponseWriter
	Req        *http.Request
	PathParams map[string]string

	// 缓存的数据
	cacheQueryValues url.Values
}

// BindJson 解析请求体中的 json 数据
func (c *Context) BindJson(val any) error {
	//if val == nil {
	//	return errors.New("web: 输入为 nil")
	//}
	if c.Req.Body == nil {
		return errors.New("web: Body 为 nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	//decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

// FormValue 解析请求体中 Form 的数据
func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return c.Req.FormValue(key), nil
}

// QueryValue 解析请求体中的 Query 数据
// 查询参数: URL 中 ? 后面的数据
// Query 和  Form 相比没有缓存
func (c *Context) QueryValue(key string) (string, error) {
	// 缓存 Query 数据 --> 避免重复 ParseQuery
	// 第一次访问时，c.cacheQueryValues 为 nil
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return "", errors.New("web: key 不存在")
	}
	return vals[0], nil
	// 用户区别不出有值但为空和没有这个参数
	//return c.Req.URL.Query().Get(key), nil
}

// PathValue 解析请求体中的 Path 数据
func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key 不存在")
	}
	return val, nil
}
