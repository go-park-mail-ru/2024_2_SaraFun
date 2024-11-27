package config

type EnvConfig struct {
	RedisUser     string `env: "REDIS_USER"`
	RedisPassword string `env: "REDIS_PASSWORD"`
	DbHost        string `env: "DB_HOST"`
	DbPort        string `env: "DB_PORT"`
	DbUser        string `env: "DB_USER"`
	DbPassword    string `env: "DB_PASSWORD"`
	DbName        string `env: "DB_NAME"`
	DbSSLMode     string `env: "DB_SSLMODE"`
}
