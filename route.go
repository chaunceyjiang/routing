package routing

import (
	"fmt"
	"net/url"
	"strings"
)

type Route struct {
	group    *RouteGroup
	name     string
	path     string // 前缀的全路径
	template string
}

func newRoute(path string, group *RouteGroup) *Route {
	// 加入group所属 路由组
	path = group.prefix + path

	name := path
	// 匹配所有
	if strings.HasSuffix(path, "*") {
		path = path[:len(path)-1] + "<:.*>"
	}

	route := &Route{group: group,
		name:     name,
		path:     path,
		template: buildURLTemplate(path)}
	group.router.routes[name] = route // 加入group所属 路由组
	return route
}

// buildURLTemplate 提取路由上面的模板信息 /a/c/<name>
func buildURLTemplate(path string) string {
	// path 子路由URL
	template, start, end := "", -1, -1
	for i := 0; i < len(path); i++ {
		if path[i] == '<' && start < 0 {
			start = i
		} else if path[i] == '>' && start >= 0 {
			name := path[start+1 : i] // 获取 <> 中的字符串

			for j := start + 1; j < i; j++ {
				// 支持 类似 /users/act-<id:\d+> (/users/act-123)
				if path[j] == ':' {
					name = path[start+1 : j]
					break
				}
			}
			// path[end+1:start]  保留URL前面的字符串
			template += path[end+1:start] + "<" + name + ">"

			// 继续寻找后面的 <>模板
			end = i
			start = -1
		}
	}

	if end < 0 {
		// 一个都没有找
		template = path
	} else if end < len(path)-1 {
		template += path[end+1:] // 保留后面的字符串
	}

	return template
}

// register 注册路由
func (r *Route) register(method string, handler ...Handler) *Route {
	hh := make([]Handler, len(handler)+len(r.group.handlers))
	// 将 路由组的 handler 和单个route 的handler 进行合并，然后给组路由下的子路由
	copy(hh, r.group.handlers)
	copy(hh[len(r.group.handlers):], handler)
	r.group.router.register(method, r.path, hh)
	return r
}

// Name 返回*Route 方便进行链式调用
func (r *Route) Name(name string) *Route {
	r.name = name
	r.group.router.routes[name] = r
	return r
}

func (r *Route) Get(handler ...Handler) *Route {
	return r.register("GET",handler...)
}

// Post Post方法 给当前的路由 增加 post方法
func (r *Route) Post(handlers ...Handler) *Route {
	return r.register("POST", handlers...)
}

// Put Put 给当前的路由 增加 Put
func (r *Route) Put(handlers ...Handler) *Route {
	return r.register("PUT", handlers...)
}

func (r *Route) Patch(handlers ...Handler) *Route {
	return r.register("PATCH", handlers...)
}

func (r *Route) Delete(handlers ...Handler) *Route {
	return r.register("DELETE", handlers...)
}

func (r *Route) Connect(handlers ...Handler) *Route {
	return r.register("CONNECT", handlers...)
}

func (r *Route) Head(handlers ...Handler) *Route {
	return r.register("HEAD", handlers...)
}

func (r *Route) Options(handlers ...Handler) *Route {
	return r.register("OPTIONS", handlers...)
}

func (r *Route) Trace(handlers ...Handler) *Route {
	return r.register("TRACE", handlers...)
}

func (r *Route) To(methods string, handler ...Handler) {
	for _, method := range strings.Split(methods, ",") {
		r.register(method, handler...)
	}
}

func (r *Route) URL(pairs ...interface{}) string {
	s := r.template
	for i := 0; i < len(pairs); i++ {
		name := fmt.Sprintf("<%v>", pairs[i])
		value := ""
		if i < len(pairs)-1 {
			value = url.QueryEscape(fmt.Sprint(pairs[i+1]))
		}
		// 将模板解析
		s = strings.Replace(s, name, value, -1)
	}
	return s
}


func combineHandlers(h1 []Handler, h2 []Handler) []Handler {
	hh := make([]Handler, len(h1)+len(h2))
	copy(hh, h1)
	copy(hh[len(h1):], h2)
	return hh
}