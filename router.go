package routing

import (
	"errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"sync"
)

/*
	Router ---> RouteGroup---> Router

	Router ----> Route -->RouteGroup

*/

type (
	Handler func(*Context) error
	Router  struct {
		RouteGroup
		routes           map[string]*Route
		pool             sync.Pool
		stores           map[string]routeStore
		maxParams        int
		notFound         []Handler
		notFoundHandlers []Handler
	}

	routeStore interface {
		Add(key string, data interface{}) int
		Get(key string, pvalues []string) (data interface{}, pnames []string)
		String() string
	}
)

func New() *Router {
	r :=&Router{
		RouteGroup:       RouteGroup{},
		routes:           make(map[string]*Route),

		stores:           make(map[string]routeStore),
		maxParams:        0,
		notFound:         nil,
	}
	r.NotFound(MethodNotAllowedHandler,NotFoundHandler)
	r.pool =  sync.Pool{New: func() interface{} {
		return &Context{pvalues:make([]string,r.maxParams),router:r}
	}}
	return r
}
// register 将 每一个路由对应的方法进行存储
func (r *Router) register(method, path string, handlers []Handler) {
	//store:=
}

func (r *Router) NotFound(handlers ...Handler) {
	r.notFound = handlers
	r.notFoundHandlers = combineHandlers(r.handlers,r.notFound)
}

func (r *Router) HandleRequest(ctx *fasthttp.RequestCtx) {
	c := r.pool.Get().(*Context)
	c.init(ctx)
	c.handlers, c.pnames = r.find(string(ctx.Method()), string(ctx.Path()), c.pvalues)

	if err := c.Next(); err != nil {
		r.handleError(c, err)
	}
	r.pool.Put(c)
}

func (r *Router) handleError(c *Context, err error) {
	var e *httpError
	if ok := errors.As(err, &e); ok {
		c.Error(e.Error(), e.StatusCode())
	} else {
		c.Error(e.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) find(method, path string, pvalues []string) (handlers []Handler, pnames []string) {
	var hh interface{}

	if store := r.stores[method]; store != nil {
		hh, pnames = store.Get(path, pvalues)
	}
	if hh != nil {
		return hh.([]Handler), pnames
	}
	return r.notFoundHandlers, pnames
}

func NotFoundHandler(ctx *Context) error {
	return NewHTTPError(http.StatusNotFound)
}

func MethodNotAllowedHandler(ctx *Context) error {
	//methods := ctx.Router().


	return nil
}