package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/mongo"
)

type Config struct {
	Mongo    *mongo.MongoConfig  `json:"mongo"`
	Logger   logger.LoggerConfig `json:"logger"`
	PassSalt string              `json:"pass_salt"`
}

func Read(path string) (c Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return c, fmt.Errorf("Failed to open the config file: %v\n", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return c, fmt.Errorf("Failed to load the config file: %v\n", err)
	}

	return c, nil
}

func MustRead(path string) Config {
	c, err := Read(path)
	if err != nil {
		panic(err)
	}
	return c
}
