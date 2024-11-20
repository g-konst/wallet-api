package internal

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
	"wallet-app/pkg/logger"
)

type Config struct {
	ListenAddr        string        `mapstructure:"LISTEN_ADDR"`
	ListenPort        int           `mapstructure:"LISTEN_PORT"`
	ServerTimeout     time.Duration `mapstructure:"SERVER_TIMEOUT"`
	ServerIdleTimeout time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT"`
	DatabaseUrl       string        `mapstructure:"DATABASE_URL"`
}

func (c *Config) ListenAndPort() string {
	return fmt.Sprintf("%s:%d", c.ListenAddr, c.ListenPort)
}

func MustLoadConfig() *Config {
	log := logger.NewLogger()

	var config Config

	v := viper.New()
	v.SetDefault("LISTEN_ADDR", "127.0.0.1")
	v.SetDefault("LISTEN_PORT", "8000")
	v.SetDefault("SERVER_TIMEOUT", "5s")
	v.SetDefault("SERVER_IDLE_TIMEOUT", "5s")
	v.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres")

	v.SetConfigType("env")
	v.SetConfigFile("config.env")
	if err := v.ReadInConfig(); err != nil {
		log.Error("cannot load config", err)
	}

	v.AutomaticEnv()
	if err := v.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
