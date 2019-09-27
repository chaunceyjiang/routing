package routing

import "sync"
/*
	Router ---> RouteGroup---> Router

	Router ----> Route -->RouteGroup

*/

type (
	Handler func(*Context) error
	Router struct {
		RouteGroup
		routes map[string]*Route
		pool sync.Pool
		stores           map[string]routeStore
		maxParams int
		notFound []Handler
		notFoundHandlers []Handler
	}

	routeStore interface {
		Add(key string,data interface{}) int
		Get(key string,pvalues []string) (data interface{},pnames []string)
		String() string
	}

)

// register 将 每一个路由对应的方法进行存储
func (r *Router) register(method, path string, handlers []Handler) {
	//store:=
}


