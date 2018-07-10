package login

import (
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/net/http"
	"github.com/linkernetworks/session"
)

func (s *LoginService) me(req *restful.Request, resp *restful.Response) {

	token := req.Request.Header.Get("Authorization")
	session, err := session.Service.Store.Get(req.Request, SessionKey)
	if err != nil {
		http.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	user := s.GetCurrentUserRestful(req)
	if user == nil {
		http.NotFound(req.Request, resp, err)
		return

	}
	resp.WriteEntity(SessionResponse{
		ID:          user.ID.Hex(),
		Token:       token,
		ExpiredAt:   session.Values["expiredAt"].(int64),
		CurrentUser: *user,
	})
}
