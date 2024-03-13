package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
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
func (c *Context) FormValue(key string) StringValue {
	err := c.Req.ParseForm()
	if err != nil {
		return StringValue{
			err: err,
		}
	}
	return StringValue{
		val: c.Req.FormValue(key),
	}
}

// QueryValue 解析请求体中的 Query 数据
// 查询参数: URL 中 ? 后面的数据
// Query 和  Form 相比没有缓存
func (c *Context) QueryValue(key string) StringValue {
	// 缓存 Query 数据 --> 避免重复 ParseQuery
	// 第一次访问时，c.cacheQueryValues 为 nil
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key 不存在"),
		}
	}
	return StringValue{
		val: vals[0],
	}
	// 用户区别不出有值但为空和没有这个参数
	//return c.Req.URL.Query().Get(key), nil
}

// PathValue 解析请求体中的 Path 数据
func (c *Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key 不存在"),
		}
	}
	return StringValue{
		val: val,
	}
}

// StringValue 结构体
type StringValue struct {
	val string
	err error
}

// ToInt64 转换为 int64
// 通过这种方式进行链式调用
// 不需要在处理输入解析每种数据都写一次不同的数据类型的方法 int64, int32...
func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
