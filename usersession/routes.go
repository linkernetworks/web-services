package usersession

import (
	restful "github.com/emicklei/go-restful"
	// "github.com/linkernetworks/vortex/src/net/http"

	"github.com/linkernetworks/web-services/serviceprovider"
	"github.com/linkernetworks/web-services/web"
)

func NewLoginService(sp *serviceprovider.Container) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/v1").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	ws.Route(ws.GET("/me").Filter(SessionAuthenticationFilter).To(RESTfulServiceHandler(sp, GetMeHandler)))
	ws.Route(ws.POST("/email/check").To(RESTfulServiceHandler(sp, CheckEmailAvailability)))
	ws.Route(ws.POST("/signup").To(RESTfulServiceHandler(sp, SignUpUserHandler)))
	ws.Route(ws.POST("/signin").To(RESTfulServiceHandler(sp, SignInUserHandler)))
	ws.Route(ws.GET("/signout").Filter(SessionAuthenticationFilter).To(RESTfulServiceHandler(sp, SignOutUserHandler)))
	return ws
}

type RESTfulContextHandler func(*web.Context)

func RESTfulServiceHandler(sp *serviceprovider.Container, handler RESTfulContextHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx := web.Context{sp, req, resp}
		handler(&ctx)
	}
}
