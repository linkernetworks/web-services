package login

import (
	"fmt"

	restful "github.com/emicklei/go-restful"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/imdario/mergo"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/webservice/login/config"
	"github.com/linkernetworks/webservice/userstorage"
)

type LoginService struct {
	passworldSalt string
	userStorage   userstorage.UserStorage
	store         sessions.Store
	web           *restful.WebService
}

func New(c *config.Config) (*LoginService, error) {

	dc := config.DefaultConfig

	if c != nil {
		if err := mergo.Merge(&dc, c, mergo.WithOverride); err != nil {
			return nil, fmt.Errorf("Generate config failed. config: [%#v] err: [%v]", c, err)
		}
	}

	logger.Debugf("Use config [%#v].", dc)

	user, err := getUserStorage(&dc.UserStore)
	if err != nil {
		return nil, fmt.Errorf("Get user storage failed. config: [%#v] err: [%v]", dc.UserStore, err)
	}

	store, err := getSessionStore(&dc.SessionStore)
	if err != nil {
		return nil, fmt.Errorf("Get session store failed. config: [%#v] err: [%v]", dc.SessionStore, err)
	}

	s := &LoginService{
		passworldSalt: dc.PassworldSalt,
		userStorage:   user,
		store:         store,
		web:           &restful.WebService{},
	}

	s.web.Path("/v1").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	s.web.Route(s.web.GET("/me").Filter(s.authenticatedFilter).To(s.me))
	s.web.Route(s.web.POST("/email/check").To(s.checkEmail))
	s.web.Route(s.web.POST("/signup").To(s.signUp))
	s.web.Route(s.web.POST("/signin").To(s.signIn))
	s.web.Route(s.web.GET("/signout").Filter(s.authenticatedFilter).To(s.signOut))
	return s, nil
}

func (s *LoginService) WebService() *restful.WebService {
	return s.web
}

func getUserStorage(c *config.StoreConfig) (userstorage.UserStorage, error) {
	switch c.Type {
	case config.MEMORY:
		logger.Debugf("Use in-memory user storage.")
		return userstorage.NewMemoryStorage(), nil
	default:
		return nil, fmt.Errorf("Type [%v] dose not support", c.Type)
	}
}

func getSessionStore(c *config.StoreConfig) (sessions.Store, error) {
	switch c.Type {
	case config.MEMORY:
		logger.Debugf("Use in-memory session store.")
		key := securecookie.GenerateRandomKey(64)
		return sessions.NewCookieStore(key), nil
	default:
		return nil, fmt.Errorf("Type [%v] dose not support", c.Type)
	}
}
