package web

import (
	"net/http"

	"github.com/linkernetworks/logger"
)

func (web *Web) SignIn(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	logger.Debugf("Get sign in request with email [%v]", email)

	if web.userAuth.Verify(email, password) {
		// TODO: return data in json format
		w.Write([]byte("ok"))
		return
	}

	// TODO: return data in json format
	http.Error(w, "failed", 400)
}
