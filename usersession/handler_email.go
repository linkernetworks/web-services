package usersession

import (
	"net/http"

	"fmt"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
	response "github.com/linkernetworks/net/http"
	oauth "github.com/linkernetworks/oauth/entity"
	"github.com/linkernetworks/validator"
	"github.com/linkernetworks/webservice/web"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type emailCheckRequest struct {
	Email string `json:"email"`
}

func CheckEmailAvailability(ctx *web.Context) {
	sp, req, resp := ctx.ServiceProvider, ctx.Request, ctx.Response

	var e emailCheckRequest
	if err := req.ReadEntity(&e); err != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, ActionResponse{
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
		resp.WriteHeaderAndEntity(http.StatusUnprocessableEntity, ActionResponse{
			Error:       true,
			Validations: validations,
			Message:     "Email is not valid",
		})
		return
	}

	session := sp.Mongo.NewSession()
	defer session.Close()

	// Check user existed
	query := bson.M{"email": email}
	existedUser := oauth.User{}
	if err := session.FindOne(oauth.UserCollectionName, query, &existedUser); err != nil {
		logger.Error(err)
		if err == mgo.ErrNotFound {
			response.NotFound(req.Request, resp.ResponseWriter, err)
			return
		}
		response.InternalServerError(req.Request, resp.ResponseWriter, err)
		return
	}
	if len(existedUser.ID) > 1 {
		msg := fmt.Sprintf("User email: %s already existed.", existedUser.Email)
		resp.WriteHeaderAndEntity(http.StatusConflict, ActionResponse{
			Error:       true,
			Validations: validations,
			Message:     msg,
		})
		return
	}

	resp.WriteEntity(ActionResponse{
		Error:   false,
		Message: "email OK",
	})
}

func loadUserByEmail(service *mongo.Service, email string) (*oauth.User, error) {
	session := service.NewSession()
	defer session.Close()

	var q = bson.M{"email": email}
	var u oauth.User
	var err = session.FindOne(oauth.UserCollectionName, q, &u)
	return &u, err
}
