package routing

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type SerializeFunc func(data interface{}) ([]byte, error)

// Context 存储每个请求的上下文环境
type Context struct {
	*fasthttp.RequestCtx

	router *Router

	Serialize SerializeFunc // 将任意类型进行序列化

	pnames  []string //参数名
	pvalues []string //参数值

	data map[string]interface{}

	handlers []Handler

	index int
}

func (c *Context) Router() *Router {
	return c.router
}

// Param 根据参数名获取值
func (c *Context) Param(name string) string {

	for i, n := range c.pnames {
		if n == name {
			return c.pvalues[i]
		}
	}
	return ""
}

func (c *Context) Get(name string) interface{} {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	return c.data[name]
}
func (c *Context) Set(name string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[name] = value
}

func (c *Context) Next() error {
	c.index++
	for n := len(c.handlers); c.index < n; c.index++ {
		if err := c.handlers[c.index](c); err != nil {
			return nil
		}
	}
	return nil
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// URL 获取该路由的URL
func (c *Context) URL(route string, pairs ...interface{}) string {
	if r := c.router.routes[route]; r != nil {
		return r.URL(pairs...)
	}
	return ""
}

func (c *Context) WriteData(data interface{}) error {
	var bytes []byte
	var err error
	if bytes, err = c.Serialize(data); err == nil {
		_, err = c.Write(bytes)
	}

	return err
}
func (c *Context) init(ctx *fasthttp.RequestCtx) {
	c.RequestCtx = ctx
	c.data = nil
	c.index = -1
	c.Serialize = Serialize
}

func Serialize(data interface{}) (bytes []byte, err error) {
	switch data.(type) {
	case []byte:
		return data.([]byte), nil
	case string:
		return []byte(data.(string)), nil
	default:
		if data != nil {
			return []byte(fmt.Sprint(data)), nil
		}
	}
	return nil, nil
}