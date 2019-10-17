package routing

import "testing"

type mockStore struct {
	*store
	data map[string]interface{}
}

func TestRouteGroup(t *testing.T) {
	router := New()
	for _,method := range
}