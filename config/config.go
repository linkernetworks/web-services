package config

import "github.com/linkernetworks/logger"

type Storage string

const (
	Memory Storage = "Memory"
	Mongo  Storage = "Mongo"
	Redis  Storage = "Redis"
)

type Config struct {
	Store  StoreConfig
	User   UserConfig
	Logger logger.LoggerConfig
}

type StoreConfig struct {
	Type Storage
}

type UserConfig struct {
	Type     Storage
	MongoURL string
}

var DefaultConfig Config = Config{
	Store: StoreConfig{
		Type: Memory,
	},
	User: UserConfig{
		Type: Memory,
	},
	Logger: logger.LoggerConfig{
		Dir:           "./logs",
		Level:         "info",
		MaxAge:        "720h",
		SuffixPattern: ".%Y%m%d",
		LinkName:      "log",
	},
}
