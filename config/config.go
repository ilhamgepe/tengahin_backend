package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

// Server config struct
type ServerConfig struct {
	AppVersion           string
	Port                 string
	Mode                 string
	JWTSecretKey         string
	JWTRefreshSecretKey  string
	TokenDuration        time.Duration
	RefreshTokenDuration time.Duration
	CookieName           string
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	CtxDefaultTimeout    time.Duration
}

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
	PgMaxConn          int
	PgMaxConnLifetime  int
	PgMaxIdleTime      time.Duration
	PgConnectTimeout   time.Duration
}

type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	Protocol      int
	DB            int
}

// Load config file from given path
func LoadConfig(path string, filename string) (*Config, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.SetConfigType("yml")
	v.AddConfigPath(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg, nil
}
