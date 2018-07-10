package config

type Type string

const (
	Memory Type = "MEMORY"
	Mongo  Type = "MONGO"
	Redis  Type = "REDIS"
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
		Type:      Memory,
		MongoURL:  "mongodb://localhost:27017/user",
		RedisHost: "localhost",
		RedisPort: 6379,
	},
	SessionStore: StoreConfig{
		Type:      Memory,
		MongoURL:  "mongodb://localhost:27017/session",
		RedisHost: "localhost",
		RedisPort: 6379,
	},
}
