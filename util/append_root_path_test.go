package util

import (
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendRootPath(t *testing.T) {
	t.Parallel()

	data := []struct {
		method string
	}{{
		"GET",
	}, {
		"POST",
	}}

	for _, d := range data {
		t.Run(d.method, func(t *testing.T) {
			t.Parallel()

			// arrange: create a webservice with one route to /old/foo
			ws := &restful.WebService{}
			ws.SetDynamicRoutes(true)
			ws.Path("/old")
			ws.Route(ws.Method(d.method).Path("/foo").To(func(*restful.Request, *restful.Response) {}))

			// action: append path from /new/old/foo
			AppendRootPath(ws, "/new")

			// assert: check route path
			require.Len(t, ws.Routes(), 1)
			assert.Equal(t, "/new/old/foo", ws.Routes()[0].Path)
		})
	}
}
