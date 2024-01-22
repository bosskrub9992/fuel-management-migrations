package config

import (
	"strings"

	"github.com/jinleejun-corp/corelib/databases"
	"github.com/jinleejun-corp/corelib/slogger"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Postgres databases.PostgresConfig
		SQLite   struct {
			FilePath string
		}
	}
	Logger slogger.Config
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config") // local
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.AutomaticEnv()
}

func New() *Config {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
