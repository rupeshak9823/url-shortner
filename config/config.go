package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	DBConfig
	RedisConfig
}

type DBConfig struct {
	DBMasterConfig
	DbLogEnable                 bool          `envconfig:"POSTGRES_LOG_MODE_ENABLE"`
	DbMaxIdleConnection         int           `envconfig:"POSTGRES_MAX_IDLE_CONNECTION"`
	DbMaxConnection             int           `envconfig:"POSTGRES_MAX_CONNECTION"`
	DbMaxConnectionLifetime     time.Duration `envconfig:"POSTGRES_MAX_CONNECTION_LIFETIME"`
	DbMaxIdleConnectionLifetime time.Duration `envconfig:"POSTGRES_MAX_IDLE_CONNECTION_LIFETIME"`
}

type DBMasterConfig struct {
	DbHost string `envconfig:"POSTGRES_HOST"`
	DbPort int    `envconfig:"POSTGRES_PORT"`
	DbName string `envconfig:"POSTGRES_DB"`
	DbUser string `envconfig:"POSTGRES_USER"`
	DbPass string `envconfig:"POSTGRES_PASS"`
}

type RedisConfig struct {
	Host string `envconfig:"REDIS_HOST"`
	Port int    `envconfig:"REDIS_PORT"`
	DB   int    `envconfig:"REDIS_DB"`
}

func GetAppConfigFromEnv() AppConfig {
	var conf AppConfig
	envconfig.MustProcess("", &conf)
	return conf
}
