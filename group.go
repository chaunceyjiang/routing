package routing


// RouteGroup 路由组 代表 一组有着同样的前缀的的路由
type RouteGroup struct {
	prefix string
	router *Router
	handlers []Handler

}