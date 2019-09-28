package routing

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	e := NewHTTPError(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, e.StatusCode())
	assert.Equal(t, http.StatusText(http.StatusNotFound), e.Error())


	e = NewHTTPError(http.StatusNotFound,"error")
	assert.Equal(t,http.StatusNotFound,e.StatusCode())
	assert.Equal(t,"error",e.Error())


	s,_ :=json.Marshal(e)
	assert.Equal(t,`{"status":404,"message":"error"}`,string(s))

}
