package session

import (
	"errors"
	"net/http"

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
		responseErrorWithStatus(resp, http.StatusBadRequest, err.Error())
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
		responseErrorWithStatus(resp, http.StatusBadRequest, err.Error())
		return
	}

	session := as.Mongo.NewSession()
	defer session.Close()

	// get user from db
	logger.Debug(as.Config.Oauth.Encryption)
	password, err := pwdutil.EncryptPasswordLegacy(form.Password)
	if err != nil {
		responseErrorWithStatus(resp, http.StatusBadRequest, err.Error())
		return
	}
	query := bson.M{"email": form.Email, "password": password}

	user := oauth.User{}
	if err := session.FindOne(oauth.UserCollectionName, query, &user); err != nil {
		if err == mgo.ErrNotFound {
			responseErrorWithStatus(resp, http.StatusForbidden, ErrInvalidUsernameOrPassword.Error())
			return
		} else {
			responseErrorWithStatus(resp, http.StatusForbidden, err.Error())
			return
		}
	}

	if user.Revoked {
		resp.WriteHeaderAndEntity(http.StatusForbidden, SignInResponse{
			Error:   true,
			Message: "This user has been revoked",
		})
		return
	}

	token, err := SignIn(req.Request, resp.ResponseWriter, &user)
	if err != nil {
		responseErrorWithStatus(resp, http.StatusInternalServerError, err.Error())
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
			responseErrorWithStatus(resp, http.StatusNotFound, err.Error())
			return
		}
		responseErrorWithStatus(resp, http.StatusInternalServerError, err.Error())
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
