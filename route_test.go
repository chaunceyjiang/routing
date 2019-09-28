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

func Test_buildURLTemplate(t *testing.T) {
	tests :=[]struct{
		path,expected string
	}{
		{"", ""},
		{"/users", "/users"},
		{"<id>", "<id>"},
		{"<id", "<id"},
		{"/users/<id>", "/users/<id>"},
		{"/users/<id:\\d+>", "/users/<id>"},
		{"/users/<:\\d+>", "/users/<>"},
		{"/users/<id>/xyz", "/users/<id>/xyz"},
		{"/users/<id:\\d+>/xyz", "/users/<id>/xyz"},
		{"/users/<id:\\d+>/<test>", "/users/<id>/<test>"},
		{"/users/<id:\\d+>/<test>/", "/users/<id>/<test>/"},
		{"/users/<id:\\d+><test>", "/users/<id><test>"},
		{"/users/<id:\\d+><test>/", "/users/<id><test>/"},
	}

	for _,test :=range tests{
		assert.Equal(t,test.expected,buildURLTemplate(test.path))
	}
}