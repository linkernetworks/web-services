package util

import (
	restful "github.com/emicklei/go-restful"
)

// AppendRootPath is an util function to append a root path to all existed routes.
// For example,, if root is '/new', a route with path '/old/foo' is replaced with '/new/old/foo'.
func AppendRootPath(ws *restful.WebService, root string) *restful.WebService {

	// Set new root
	ws.Path(root)

	// Copy all existed routes
	oldRoutes := ws.Routes()

	// Reset all existed routes
	for _, oldRoute := range oldRoutes {

		if err := ws.RemoveRoute(oldRoute.Path, oldRoute.Method); err != nil {
			panic(err)
		}

		newRoute := ws.Method(oldRoute.Method).
			Consumes(oldRoute.Consumes...).
			Produces(oldRoute.Produces...).
			Path(oldRoute.Path).
			To(oldRoute.Function).
			Doc(oldRoute.Doc).
			Notes(oldRoute.Notes).
			Operation(oldRoute.Operation)

		for _, f := range oldRoute.Filters {
			newRoute.Filter(f)
		}

		for _, f := range oldRoute.If {
			newRoute.If(f)
		}

		for _, p := range oldRoute.ParameterDocs {
			newRoute.Param(p)
		}

		for _, r := range oldRoute.ResponseErrors {
			newRoute.Returns(r.Code, r.Message, r.Model)
		}

		for k, v := range oldRoute.Metadata {
			newRoute.Metadata(k, v)
		}

		ws.Route(newRoute)
	}

	return ws
}
