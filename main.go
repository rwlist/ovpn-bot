package main

import (
	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/updates"
	"log"

	"github.com/rwlist/ovpn-bot/conf"
)

func main() {
	cfg, err := conf.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	bot := telegram.NewBot(cfg.BotToken)

	ch, err := updates.StartPolling(bot, telegram.GetUpdatesRequest{
		Offset:         0,
		Limit:          50,
		Timeout:        60,
	})
	if err != nil {
		log.Fatal(err)
	}

	for upd := range ch {

	}
}