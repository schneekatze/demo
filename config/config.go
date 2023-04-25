package config

import (
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var c config

type config struct {
	LogLevel string `split_words:"true" default:"debug"`
	Listen   string `default:"0.0.0.0:8000"`
	DBConfig struct {
		Name     string `envconfig:"DATABASE_NAME"`
		Host     string `envconfig:"DATABASE_HOST"`
		Port     string `envconfig:"DATABASE_PORT"`
		User     string `envconfig:"DATABASE_USER"`
		Password string `envconfig:"DATABASE_PASSWORD"`
	}
}

func Cfg() config {
	return c
}

func init() {
	c = new()

	lvl, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(lvl)

	formattedConfig, _ := json.MarshalIndent(c, "", "    ")
	log.Debug("Configuration: ", string(formattedConfig))
}

func new() config {
	var c config

	_ = godotenv.Load(".env.local")
	_ = godotenv.Load()

	err := envconfig.Process("CHLNG", &c)

	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
