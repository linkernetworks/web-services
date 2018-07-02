package web

import (
	"net/http"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/web-services/entity"
)

func (web *Web) SignUp(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	user := entity.User{
		Email:    email,
		Password: password,
	}
	logger.Debugf("Get sign up request with user [%#v]", user)

	err := web.userAuth.Register(user)
	if err != nil {
		logger.Debugf("Register failed. user: [%#v], err: [%v]", user, err)
		// TODO: return data in json format
		http.Error(w, "failed", 400)
	}

	// TODO: return data in json format
	w.Write([]byte("ok"))
}
