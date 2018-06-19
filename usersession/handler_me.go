package usersession

import (
	"bitbucket.org/linkernetworks/aurora/src/acl"
	"github.com/linkernetworks/net/http"
	"bitbucket.org/linkernetworks/aurora/src/web"
	"github.com/linkernetworks/session"

	"gopkg.in/mgo.v2"
)

func GetMeHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	token := req.Request.Header.Get("Authorization")
	session, err := session.Service.Store.Get(req.Request, SessionKey)
	if err != nil {
		http.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	ses := as.Mongo.NewSession()
	defer ses.Close()

	user, err := acl.GetCurrentUserRestful(ses, req)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.NotFound(req.Request, resp, err)
			return
		}
		http.InternalServerError(req.Request, resp, err)
		return
	}
	resp.WriteEntity(SessionResponse{
		ID:          user.ID.Hex(),
		Token:       token,
		ExpiredAt:   session.Values["expiredAt"].(int64),
		CurrentUser: *user,
	})
}
