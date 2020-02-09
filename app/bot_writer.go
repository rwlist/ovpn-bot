package app

import (
	"github.com/petuhovskiy/telegram"
	"sync"
)

type botWriter struct {
	bot *telegram.Bot
	chatID int

	buf []byte
	m sync.Mutex
}

func NewBotWriter(bot *telegram.Bot, chatID int) *botWriter {
	return &botWriter{
		bot:    bot,
		chatID: chatID,
	}
}

func (w *botWriter) Write(b []byte) (n int, err error) {
	defer w.Flush()

	w.m.Lock()
	defer w.m.Unlock()

	w.buf = append(w.buf, b...)
	return len(b), nil
}

func (w *botWriter) Flush() {
	w.m.Lock()
	defer w.m.Unlock()

	if len(w.buf) == 0 {
		return
	}

	_, _ = w.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID:                str(w.chatID),
		Text:                  string(w.buf),
	})

	w.buf = nil
}
