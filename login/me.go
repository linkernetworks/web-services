package usersession

import (
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/net/http"
	"github.com/linkernetworks/session"
	"github.com/linkernetworks/webservice/login/entity"

	"gopkg.in/mgo.v2"
)

func (s *LoginService) me(req *restful.Request, resp *restful.Response) {

	token := req.Request.Header.Get("Authorization")
	session, err := session.Service.Store.Get(req.Request, SessionKey)
	if err != nil {
		http.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	ses := s.mongo.NewSession()
	defer ses.Close()

	user, err := s.GetCurrentUserRestful(ses, req)
	if err != nil {
		if err == mgo.ErrNotFound {
			http.NotFound(req.Request, resp, err)
			return
		}
		http.InternalServerError(req.Request, resp, err)
		return
	}
	resp.WriteEntity(entity.SessionResponse{
		ID:          user.ID.Hex(),
		Token:       token,
		ExpiredAt:   session.Values["expiredAt"].(int64),
		CurrentUser: *user,
	})
}
