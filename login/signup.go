package login

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"
	response "github.com/linkernetworks/net/http"
	"github.com/linkernetworks/validator"
	"github.com/linkernetworks/webservice/login/entity"
	"github.com/linkernetworks/webservice/pwdutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (s *LoginService) signUp(req *restful.Request, resp *restful.Response) {

	user := entity.User{}
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

	session := s.mongo.NewSession()
	defer session.Close()

	// Check user email validate first
	emailValidate, err := validator.ValidateEmail(user.Email)
	if err != nil {
		validations["email"] = emailValidate
	}
	// Then Check user existed
	query := bson.M{"email": user.Email}
	existedUser := entity.User{}
	if err := session.FindOne(entity.UserCollectionName, query, &existedUser); err != nil {
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
		resp.WriteHeaderAndEntity(http.StatusUnprocessableEntity, entity.ActionResponse{
			Error:       true,
			Validations: validations,
			Message:     "Input data is not valid",
		})
		return
	}

	user.ID = bson.NewObjectId()
	user.Password, err = pwdutil.EncryptPasswordLegacy(user.Password, s.passworldSalt)
	if err != nil {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	user.CreatedAt = time.Now().Unix()
	user.Roles = []string{"user"}
	user.Revoked = false
	user.JobPriority = 3000

	if err := session.Insert(entity.UserCollectionName, &user); err != nil {
		logger.Error(err)
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	resp.WriteEntity(entity.ActionResponse{
		Error:   false,
		Message: "Sign up success",
	})

}
