package config

type Type string

const (
	MEMORY Type = "MEMORY"
	MONGO  Type = "MONGO"
	REDIS  Type = "REDIS"
)

type Config struct {
	PassworldSalt string
	UserStore     StoreConfig
	SessionStore  StoreConfig
}

type StoreConfig struct {
	Type      Type
	MongoURL  string
	RedisHost string
	RedisPort int
}

var DefaultConfig Config = Config{
	PassworldSalt: "",
	UserStore: StoreConfig{
		Type:      MEMORY,
		MongoURL:  "mongodb://localhost:27017/user",
		RedisHost: "localhost",
		RedisPort: 6379,
	},
	SessionStore: StoreConfig{
		Type:      MEMORY,
		MongoURL:  "mongodb://localhost:27017/session",
		RedisHost: "localhost",
		RedisPort: 6379,
	},
}
