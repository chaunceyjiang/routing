package routing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteNew(t *testing.T) {
	router :=New()
	group := newRouteGroup("/admin",router,nil)

	r1 := newRoute("/users",group)

	assert.Equal(t,r1.name ,"/admin/users")
	assert.Equal(t, "/admin/users", r1.name, "route.name =")
	assert.Equal(t, "/admin/users", r1.path, "route.path =")
	assert.Equal(t, "/admin/users", r1.template, "route.template =")
	_, exists := router.routes[r1.name]
	assert.True(t, exists, "router.routes[name] is ")

	r2 := newRoute("/users/<id:\\d+>/*", group)
	assert.Equal(t, "/admin/users/<id:\\d+>/*", r2.name, "route.name =")
	assert.Equal(t, "/admin/users/<id:\\d+>/<:.*>", r2.path, "route.path =")
	assert.Equal(t, "/admin/users/<id>/<>", r2.template, "route.template =")
	_, exists = router.routes[r2.name]
	assert.True(t, exists, "router.routes[name] is ")
}