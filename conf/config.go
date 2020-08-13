package conf

import (
	"github.com/caarlos0/env/v6"
)

type Struct struct {
	AdminIDs []int  `env:"ADMIN_TELEGRAM_ID" envSeparator:","` // Comma-separated list of the bot admins' Telegram IDs
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
