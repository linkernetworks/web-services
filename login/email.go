package login

import (
	"net/http"

	"fmt"

	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/mongo"
	response "github.com/linkernetworks/net/http"
	"github.com/linkernetworks/validator"
	"github.com/linkernetworks/webservice/login/entity"
	"gopkg.in/mgo.v2/bson"
)

type emailCheckRequest struct {
	Email string `json:"email"`
}

func (s *LoginService) checkEmail(req *restful.Request, resp *restful.Response) {

	var e emailCheckRequest
	if err := req.ReadEntity(&e); err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, entity.ActionResponse{
			Error:   true,
			Message: "Failed to process entity",
		})
		return
	}

	email := e.Email

	validations := validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(email)
	if err != nil {
		validations["email"] = emailValidate
	}

	if validations.HasError() {
		resp.WriteHeaderAndEntity(http.StatusUnprocessableEntity, entity.ActionResponse{
			Error:       true,
			Validations: validations,
			Message:     "Email is not valid",
		})
		return
	}

	existedUser := s.userStorage.FindByEmail(email)
	if existedUser == nil {
		response.NotFound(req.Request, resp.ResponseWriter, err)
		return
	}
	if len(existedUser.ID) > 1 {
		msg := fmt.Sprintf("User email: %s already existed.", existedUser.Email)
		resp.WriteHeaderAndEntity(http.StatusConflict, entity.ActionResponse{
			Error:       true,
			Validations: validations,
			Message:     msg,
		})
		return
	}

	resp.WriteEntity(entity.ActionResponse{
		Error:   false,
		Message: "email OK",
	})
}

func loadUserByEmail(service *mongo.Service, email string) (*entity.User, error) {
	session := service.NewSession()
	defer session.Close()

	var q = bson.M{"email": email}
	var u entity.User
	var err = session.FindOne(entity.UserCollectionName, q, &u)
	return &u, err
}
