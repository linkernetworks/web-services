package usersession

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	response "github.com/linkernetworks/net/http"
	"github.com/linkernetworks/session"
	"github.com/linkernetworks/webservice/login/entity"
)

func (s *LoginService) signOut(req *restful.Request, resp *restful.Response) {

	sess, err := session.Service.Store.Get(req.Request, SessionKey)
	if err != nil {
		logger.Errorf("Redis get auth token failed: %v", err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	if err := session.Service.Store.Delete(req.Request, resp.ResponseWriter, sess); err != nil {
		logger.Errorf("Failed to delete token: %v", err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteHeaderAndEntity(http.StatusOK, entity.SignInResponse{
		Error:     false,
		Message:   "Logout success",
		SignInUrl: "/signin",
	})
}
