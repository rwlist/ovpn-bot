package conf

import (
	"github.com/caarlos0/env/v6"
)

type Struct struct {
	AdminID int `env:"ADMIN_TELEGRAM_ID"`
	BotToken string `env:"BOT_TOKEN"`
}

func ParseEnv() (*Struct, error) {
	cfg := Struct{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}