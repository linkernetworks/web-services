package login

import (
	"errors"
	"fmt"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/net/http"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/oauth/util"
	"github.com/linkernetworks/oauth/validator"
	"github.com/linkernetworks/webservice/login/entity"
	"github.com/linkernetworks/webservice/pwdutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Login failed. Incorrect username or password.")
)

func (s *LoginService) signIn(req *restful.Request, resp *restful.Response) {

	form := oauth.User{}
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

	session := s.mongo.NewSession()
	defer session.Close()

	// get user from db
	password, err := pwdutil.EncryptPasswordLegacy(form.Password, s.passworldSalt)

	if err != nil {
		http.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	query := bson.M{"email": form.Email, "password": password}

	user := oauth.User{}
	if err := session.FindOne(oauth.UserCollectionName, query, &user); err != nil {
		if err == mgo.ErrNotFound {
			http.Forbidden(req.Request, resp.ResponseWriter, ErrInvalidUsernameOrPassword)
			return
		} else {
			http.Forbidden(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	if user.Revoked {
		http.Forbidden(req.Request, resp, fmt.Errorf("This user has been revoked."))
		return
	}

	token, err := s.signInSession(req.Request, resp.ResponseWriter, &user)
	if err != nil {
		http.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	// Update last login timestamp & token
	user.LastLoggedInAt = util.GetCurrentTimestamp()
	user.AccessToken = token.String()

	query = bson.M{"_id": user.ID}
	modifier := bson.M{"$set": user}
	if err := session.C(oauth.UserCollectionName).Update(query, modifier); err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			http.NotFound(req.Request, resp.ResponseWriter, err)
			return
		}
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
