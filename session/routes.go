package session

import (
	"bitbucket.org/linkernetworks/aurora/src/aurora/serviceprovider"
	"bitbucket.org/linkernetworks/aurora/src/net/http"
	restful "github.com/emicklei/go-restful"
)

func NewLoginService(sp *serviceprovider.Container) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/v1").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	ws.Route(ws.GET("/me").Filter(SessionAuthenticationFilter).To(http.RESTfulServiceHandler(sp, GetMeHandler)))
	ws.Route(ws.POST("/email/check").To(http.RESTfulServiceHandler(sp, CheckEmailAvailability)))
	ws.Route(ws.POST("/signup").To(http.RESTfulServiceHandler(sp, SignUpUserHandler)))
	ws.Route(ws.POST("/signin").To(http.RESTfulServiceHandler(sp, SignInUserHandler)))
	ws.Route(ws.GET("/signout").Filter(SessionAuthenticationFilter).To(http.RESTfulServiceHandler(sp, SignOutUserHandler)))
	return ws
}
