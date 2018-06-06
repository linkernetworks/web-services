package session

import (
	"errors"
	"fmt"

	"bitbucket.org/linkernetworks/aurora/src/net/http/response"
	"bitbucket.org/linkernetworks/aurora/src/pwdutil"
	"bitbucket.org/linkernetworks/aurora/src/web"
	"github.com/linkernetworks/logger"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/oauth/util"
	"github.com/linkernetworks/oauth/validator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Login failed. Incorrect username or password.")
)

func SignInUserHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	form := oauth.User{}
	if err := req.ReadEntity(&form); err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
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
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	session := as.Mongo.NewSession()
	defer session.Close()

	// get user from db
	logger.Debug(as.Config.Oauth.Encryption)
	password, err := pwdutil.EncryptPasswordLegacy(form.Password)
	if err != nil {
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}
	query := bson.M{"email": form.Email, "password": password}

	user := oauth.User{}
	if err := session.FindOne(oauth.UserCollectionName, query, &user); err != nil {
		if err == mgo.ErrNotFound {
			response.Forbidden(req.Request, resp.ResponseWriter, ErrInvalidUsernameOrPassword)
			return
		} else {
			response.Forbidden(req.Request, resp.ResponseWriter, err)
			return
		}
	}

	if user.Revoked {
		response.Forbidden(req.Request, resp, fmt.Errorf("This user has been revoked."))
		return
	}

	token, err := SignIn(req.Request, resp.ResponseWriter, &user)
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
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
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		}
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}

	resp.WriteEntity(SignInResponse{
		Error:   false,
		Message: "Login success",
		Session: SessionResponse{
			ID:    user.ID.Hex(),
			Token: token.String(),
		},
	})
}
