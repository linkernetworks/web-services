package usersession

import (
	"net/http"

	response "github.com/linkernetworks/net/http"
	"bitbucket.org/linkernetworks/aurora/src/pwdutil"
	"bitbucket.org/linkernetworks/aurora/src/validator"
	"bitbucket.org/linkernetworks/aurora/src/web"
	"github.com/linkernetworks/logger"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/oauth/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func SignUpUserHandler(ctx *web.Context) {
	as, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	user := oauth.User{}
	if err := req.ReadEntity(&user); err != nil {
		logger.Error(err)
		response.BadRequest(req.Request, resp.ResponseWriter, err)
		return
	}

	validations := validator.ValidationMap{}
	firstNameValidate, err := validator.ValidateRequiredStrField("firstName", user.FirstName)
	if err != nil {
		validations["firstName"] = firstNameValidate
	}
	lastNameValidate, err := validator.ValidateRequiredStrField("lastName", user.LastName)
	if err != nil {
		validations["lastName"] = lastNameValidate
	}

	session := as.Mongo.NewSession()
	defer session.Close()

	// Check user email validate first
	emailValidate, err := validator.ValidateEmail(user.Email)
	if err != nil {
		validations["email"] = emailValidate
	}
	// Then Check user existed
	query := bson.M{"email": user.Email}
	existedUser := oauth.User{}
	if err := session.FindOne(oauth.UserCollectionName, query, &existedUser); err != nil {
		if err.Error() != mgo.ErrNotFound.Error() {
			logger.Error(err)
			response.InternalServerError(req.Request, resp.ResponseWriter, err)
			return
		}
	}
	if len(existedUser.ID) > 1 {
		emailValidate.Field = "email"
		emailValidate.Error = true
		emailValidate.Message = "The email is already used."
		validations["email"] = emailValidate
	}

	passworkValidate, err := validator.ValidatePassword(user.Password)
	if err != nil {
		validations["password"] = passworkValidate
	}
	if validations.HasError() {
		resp.WriteHeaderAndEntity(http.StatusUnprocessableEntity, ActionResponse{
			Error:       true,
			Validations: validations,
			Message:     "Input data is not valid",
		})
		return
	}

	user.ID = bson.NewObjectId()
	user.Password, err = pwdutil.EncryptPasswordLegacy(user.Password)
	if err != nil {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	user.CreatedAt = util.GetCurrentTimestamp()
	user.Roles = []string{"user"}
	user.Revoked = false
	user.JobPriority = 3000

	if err := session.Insert(oauth.UserCollectionName, &user); err != nil {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "Sign up success",
	})

}
