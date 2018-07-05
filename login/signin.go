package login

import (
	"errors"
	"fmt"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/net/http"
	"github.com/linkernetworks/validator"
	"github.com/linkernetworks/webservice/login/entity"
	"github.com/linkernetworks/webservice/pwdutil"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Login failed. Incorrect username or password.")
)

func (s *LoginService) signIn(req *restful.Request, resp *restful.Response) {

	form := entity.User{}
	if err := req.ReadEntity(&form); err != nil {
		http.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	var validations = validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(form.Email)
	if err != nil {
		validations["email"] = emailValidate
	}
	passworkValidate, err := validator.ValidatePassword(form.Password)
	if err != nil {
		validations["password"] = passworkValidate
	}
	if validations.HasError() {
		logger.Error(err)
		http.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	// get user from db
	password, err := pwdutil.EncryptPasswordLegacy(form.Password, s.passworldSalt)

	if err != nil {
		http.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	user := s.userStorage.FindByPassword(form.Email, password)
	if user == nil {
		http.Forbidden(req.Request, resp.ResponseWriter, ErrInvalidUsernameOrPassword)
		return
	}

	if user.Revoked {
		http.Forbidden(req.Request, resp, fmt.Errorf("This user has been revoked."))
		return
	}

	token, err := s.signInSession(req.Request, resp.ResponseWriter, user)
	if err != nil {
		http.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	// Update last login timestamp & token
	user.LastLoggedInAt = time.Now().Unix()
	user.AccessToken = token.String()

	err = s.userStorage.Save(user)
	if err != nil {
		http.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(entity.SignInResponse{
		Error:   false,
		Message: "Login success",
		Session: entity.SessionResponse{
			ID:    user.ID.Hex(),
			Token: token.String(),
		},
	})
}
