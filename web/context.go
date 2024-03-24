package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Resp           http.ResponseWriter
	Req            *http.Request
	RespStatusCode int
	RespData       []byte

	PathParams map[string]string

	// 缓存的数据
	cacheQueryValues url.Values

	MatchedRoute string

	// 页面渲染引擎
	tplEngine TemplateEngine

	// 用户可以自由决定在这里存储什么，主要用于解决在不同 Middleware 之间数据传递的问题
	// 但是 UserValues 在初始状态的时候总是 nil，需要手动初始化
	UserValues map[string]any
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.Resp, c.Req, url, http.StatusFound)
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

// SetCookie 设置 Cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

// RespJSONOK code:200
func (c *Context) RespJSONOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

// RespJSON 处理输出,返回 JSON 数据
func (c *Context) RespJSON(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	//c.Resp.WriteHeader(code)
	//_, err = c.Resp.Write(bs)
	//return err
	c.RespStatusCode = code
	c.RespData = bs
	return err
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

func (s StringValue) ToUInt64() (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseUint(s.val, 10, 64)
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

// 不能用泛型
// func (s StringValue) To[T any]() (T, error) {
//
// }

// Render 渲染模板
func (c *Context) Render(tplName string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tplName, data)
	if err != nil {
		c.RespStatusCode = 500
		return err
	}
	c.RespStatusCode = 200
	return nil
}
