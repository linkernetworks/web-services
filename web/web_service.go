package web

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/imdario/mergo"
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/web-services/authenticator"
	"github.com/linkernetworks/web-services/config"
	"github.com/linkernetworks/web-services/store"
)

type Web struct {
	config   config.Config
	store    sessions.Store
	userAuth authenticator.Authenticator
}

func New(c *config.Config) (*Web, error) {

	var err error

	s := &Web{
		config: config.DefaultConfig,
	}

	if err = mergo.Merge(&s.config, c, mergo.WithOverride); err != nil {
		return nil, fmt.Errorf("Merge config failed. err: [%v]", err)
	}

	logger.Setup(s.config.Logger)

	logger.Debugf("Using config: [%#v]", s.config)

	s.store, err = createStore(&s.config.Store)
	if err != nil {
		return nil, fmt.Errorf("Create store failed. config: [%#v] err: [%v]", s.config.Store, err)
	}

	s.userAuth, err = createUserAuth(&s.config.User)
	if err != nil {
		return nil, fmt.Errorf("Create user authenticator failed. config: [%#v] err: [%v]", s.config.User, err)
	}

	return s, nil
}

func createStore(c *config.StoreConfig) (sessions.Store, error) {
	switch c.Type {
	case config.Memory:
		logger.Infoln("Use in-memory session store")
		return store.NewMemoryStore(), nil
	default:
		return nil, fmt.Errorf("Store [%v] dose not support", c.Type)
	}
}

func createUserAuth(c *config.UserConfig) (authenticator.Authenticator, error) {
	switch c.Type {
	case config.Memory:
		logger.Infoln("Use in-memory user authenticator")
		return authenticator.NewMemory()
	case config.Mongo:
		logger.Infoln("Use mongo user authenticator")
		return authenticator.NewMongo(c.MongoURL)
	default:
		return nil, fmt.Errorf("Authenticator [%v] dose not support", c.Type)
	}
}
