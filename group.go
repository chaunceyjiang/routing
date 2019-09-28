package routing


// RouteGroup 路由组 代表 一组有着同样的前缀的的路由
type RouteGroup struct {
	prefix string
	router *Router
	handlers []Handler

}


// 给一个路由创建一个组 用着共同的前缀
func newRouteGroup(prefix string, router *Router, handlers []Handler) *RouteGroup {
	// 在router 路由下面创建路由组
	return &RouteGroup{
		prefix:   prefix,
		router:   router,  // 上级路由
		handlers: handlers,
	}
}

func (g *RouteGroup) Get(path string, handlers ...Handler) *Route {
	return newRoute(path,g).Get(handlers...)
}

func (r *RouteGroup) Post(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Post(handlers...)
}

func (r *RouteGroup) Put(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Put(handlers...)
}

func (r *RouteGroup) Patch(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Patch(handlers...)
}

func (r *RouteGroup) Delete(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Delete(handlers...)
}

func (r *RouteGroup) Connect(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Connect(handlers...)
}

func (r *RouteGroup) Head(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Head(handlers...)
}


func (r *RouteGroup) Options(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Options(handlers...)
}


func (r *RouteGroup) Trace(path string, handlers ...Handler) *Route {
	return newRoute(path, r).Trace(handlers...)
}