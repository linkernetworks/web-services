package session

import (
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/web"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/session"
)

func SignOutUserHandler(ctx *web.Context) {
	req, resp := ctx.Request, ctx.Response

	sess, err := session.Service.Store.Get(req.Request, SessionKey)
	if err != nil {
		logger.Errorf("Redis get auth token failed: %v", err)
		responseErrorWithStatus(resp, http.StatusInternalServerError, err.Error())
		return
	}

	if err := session.Service.Store.Delete(req.Request, resp.ResponseWriter, sess); err != nil {
		logger.Errorf("Failed to delete token: %v", err)
		responseErrorWithStatus(resp, http.StatusInternalServerError, err.Error())
		return
	}

	resp.WriteHeaderAndEntity(http.StatusOK, SignInResponse{
		Error:     false,
		Message:   "Logout success",
		SignInUrl: "/signin",
	})
}
