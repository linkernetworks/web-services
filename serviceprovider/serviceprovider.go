package serviceprovider

import (
	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
	"github.com/linkernetworks/web-services/config"
)

type Container struct {
	Config config.Config
	Mongo  *mongo.Service
}

type ServiceDiscoverResponse struct {
	// Container map[string]Service `json:"services"`
}

type Service interface{}

type SmartTrackerService struct {
	// Redis     serviceconfig.ServiceConfig `json:"redis"`
	// Gearman   serviceconfig.ServiceConfig `json:"gearman"`
	// Memcached serviceconfig.ServiceConfig `json:"memcached"`
}

func New(cf config.Config) *Container {
	// setup logger configuration
	logger.Setup(cf.Logger)

	logger.Infof("Connecting to mongodb: %s", cf.Mongo.Url)
	mongo := mongo.New(cf.Mongo.Url)

	sp := &Container{
		Config: cf,
		Mongo:  mongo,
	}

	return sp
}

func NewContainer(configPath string) *Container {
	cf := config.MustRead(configPath)
	return New(cf)
}
