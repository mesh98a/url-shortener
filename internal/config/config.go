package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local" env-required:"true" `
	MYSQL      `yaml:"mysql" `
	HTTPServer `yaml:"http_server"`
}

type MYSQL struct {
	Host     string `yaml:"host"  env-default:"127.0.0.1"`
	Port     string `yaml:"port"  env-default:"3306"`
	User     string `yaml:"user"  env-required:"true"`
	Password string `yaml:"password"`
}

type HTTPServer struct {
	Addr        string        `yaml:"addr" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() Config {
	conf := os.Getenv("CONFIG_PATH")
	if conf == "" {
		slog.Error("config file not exist")
	}

	if _, err := os.Stat(conf); err != nil {
		slog.Error("config file not exist")
	}
	var config Config

	if err := cleanenv.ReadConfig(conf, &config); err != nil {
		slog.Error("cannot reasd congif:", err)
	}
	return config
}
