package app

import (
	"github.com/petuhovskiy/telegram"
	"sync"
	"time"
)

const flushTimeout = 2 * time.Second

type botWriterState uint

const (
	stateNoop botWriterState = iota
	stateAwaitingFlush
)

type botWriter struct {
	bot    *telegram.Bot
	chatID int

	buf []byte
	m   sync.Mutex
	state botWriterState
}

func NewBotWriter(bot *telegram.Bot, chatID int) *botWriter {
	return &botWriter{
		bot:    bot,
		chatID: chatID,
		state:  stateNoop,
	}
}

func (w *botWriter) Write(b []byte) (n int, err error) {
	defer w.DeferFlush()

	w.m.Lock()
	defer w.m.Unlock()

	w.buf = append(w.buf, b...)
	return len(b), nil
}

func (w *botWriter) Flush() {
	w.m.Lock()
	defer w.m.Unlock()

	w.state = stateNoop

	if len(w.buf) == 0 {
		return
	}

	_, _ = w.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID: str(w.chatID),
		Text:   string(w.buf),
	})

	w.buf = nil
}

func (w *botWriter) DeferFlush() {
	w.m.Lock()
	defer w.m.Unlock()

	if w.state != stateNoop {
		return
	}
	w.state = stateAwaitingFlush

	time.AfterFunc(flushTimeout, func() {
		w.Flush()
	})
}
