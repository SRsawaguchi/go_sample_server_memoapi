package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var conf *Config = nil

type PostgresConnectionString string

// config represents config of this service.
type Config struct {
	appPort   int
	appHost   string
	RdbConfig *RdbConfig
}

// rdbConfig represents config of relational database.
type RdbConfig struct {
	port     int
	host     string
	user     string
	password string
	database string
}

// Load loads config.
func Load() Config {
	if conf != nil {
		return *conf
	}

	viper.SetEnvPrefix("MEMO")
	viper.AutomaticEnv()
	c := Config{
		appPort: viper.GetInt("APP_PORT"), // 実際の環境変数名はPrefixがついた`MEMO_APP_PORT`になる。
		appHost: viper.GetString("APP_HOST"),
		RdbConfig: &RdbConfig{
			port:     viper.GetInt("DB_PORT"),
			host:     viper.GetString("DB_HOST"),
			user:     viper.GetString("DB_USER"),
			password: viper.GetString("DB_PASSWORD"),
			database: viper.GetString("DB_NAME"),
		},
	}

	return c
}

// AppPort return port number of this service.
func (c *Config) AppPort() int {
	return c.appPort
}

// AppHost return app host if this service.
func (c *Config) AppHost() string {
	return c.appHost
}

// ConnectionString return string for connecting database.
func (r *RdbConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		r.user, r.password, r.host, r.port, r.database,
	)
}
