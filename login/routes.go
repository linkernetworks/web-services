package usersession

import (
	restful "github.com/emicklei/go-restful"
	"github.com/gorilla/sessions"

	"github.com/linkernetworks/mongo"
)

type Config struct {
	PassworldSalt string
	Mongo         *mongo.MongoConfig
}

type LoginService struct {
	passworldSalt string
	mongo         *mongo.Service
	store         sessions.Store
	restful.WebService
}

func New(c *Config) *LoginService {

	s := &LoginService{}
	s.Path("/v1").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	s.Route(s.GET("/me").Filter(s.authenticatedFilter).To(s.me))
	s.Route(s.POST("/email/check").To(s.checkEmail))
	s.Route(s.POST("/signup").To(s.signUp))
	s.Route(s.POST("/signin").To(s.signIn))
	s.Route(s.GET("/signout").Filter(s.authenticatedFilter).To(s.signOut))
	return s
}
