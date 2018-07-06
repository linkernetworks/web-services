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
	restful.WebService
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
	}

	s.Path("/v1").Consumes(restful.MIME_JSON, restful.MIME_JSON).Produces(restful.MIME_JSON, restful.MIME_JSON)
	s.Route(s.GET("/me").Filter(s.authenticatedFilter).To(s.me))
	s.Route(s.POST("/email/check").To(s.checkEmail))
	s.Route(s.POST("/signup").To(s.signUp))
	s.Route(s.POST("/signin").To(s.signIn))
	s.Route(s.GET("/signout").Filter(s.authenticatedFilter).To(s.signOut))
	return s, nil
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
