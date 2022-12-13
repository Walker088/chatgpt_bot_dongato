package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	BotToken            string `env:"BOT_TOKEN"`
	OpenaiSessionCookie string `env:"OPENAI_SESSION_COOKIE"`
}

func init() {
	err := godotenv.Load(".project.env")
	if err != nil {
		log.Fatal("Error loading .project.env file")
	}
}

func GetAppConfig() *AppConfig {
	cfg := AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}
	return &cfg
}