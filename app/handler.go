package app

import (
	"fmt"
	"strings"

	"github.com/petuhovskiy/telegram"
)

type Handler struct {
	bot   *telegram.Bot
	logic *Logic

	adminIDs []int
}

func NewHandler(bot *telegram.Bot, logic *Logic, adminIDs []int) *Handler {
	return &Handler{
		bot:      bot,
		logic:    logic,
		adminIDs: adminIDs,
	}
}

func (h *Handler) Handle(upd telegram.Update) {
	if upd.Message == nil {
		return
	}

	msg := upd.Message

	isAdmin := false
	for _, adminID := range h.adminIDs {
		if msg.From.ID == adminID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return
	}

	h.handleMessage(msg)
}

func (h *Handler) sendMessage(chatID int, text string) {
	_, _ = h.bot.SendMessage(&telegram.SendMessageRequest{
		ChatID: str(chatID),
		Text:   text,
	})
}

func (h *Handler) handleMessage(msg *telegram.Message) {
	text := msg.Text
	if !strings.HasPrefix(text, "/") {
		return
	}

	cmds := strings.Split(text, " ")
	h.handleCommand(msg.Chat.ID, cmds)
}

func (h *Handler) handleCommand(chatID int, cmds []string) {
	if len(cmds) == 0 {
		return
	}

	cmd := cmds[0]
	switch cmd {
	case "/init":
		addr := ""
		if len(cmds) > 1 {
			addr = cmds[1]
		}
		h.commandInit(chatID, addr)

	case "/status":
		h.commandStatus(chatID)

	case "/generate":
		profileName := ""
		if len(cmds) > 1 {
			profileName = cmds[1]
		}
		h.commandGenerate(chatID, profileName)

	case "/remove":
		h.commandRemove(chatID)

	default:
		h.commandNotFound(chatID)
	}
}

func (h *Handler) commandInit(chatID int, addr string) {
	go func() {
		w := NewBotWriter(h.bot, chatID)

		err := h.logic.CommandInit(w, addr)
		if err != nil {
			text := fmt.Sprintf("Error while init:\n\n%s", err)
			h.sendMessage(chatID, text)
			return
		}
	}()
}

func (h *Handler) commandStatus(chatID int) {
	res, err := h.logic.CommandStatus()
	if err != nil {
		text := fmt.Sprintf("Error while status:\n\n%s", err)
		h.sendMessage(chatID, text)
		return
	}

	h.sendMessage(chatID, res)
}

func (h *Handler) commandGenerate(chatID int, profileName string) {
	if profileName == "" {
		h.sendMessage(chatID, "Please provide profileName")
		return
	}

	w := NewBotWriter(h.bot, chatID)

	res, err := h.logic.CommandGenerate(w, profileName)
	if err != nil {
		text := fmt.Sprintf("Error while generate:\n\n%s", err)
		h.sendMessage(chatID, text)
		return
	}

	_, _ = h.bot.SendDocument(&telegram.SendDocumentRequest{
		ChatID:   str(chatID),
		Document: NewBytesUploader(profileName+".ovpn", res),
	})
}

func (h *Handler) commandRemove(chatID int) {
	go func() {
		w := NewBotWriter(h.bot, chatID)
		h.logic.CommandRemove(w)
	}()
}

func (h *Handler) commandNotFound(chatID int) {
	h.commandHelp(chatID)
}

func (h *Handler) commandHelp(chatID int) {
	str := `Need some help?

/init tcp://a.b:80	initialize everything
/status				displays "docker ps"
/generate <name> 	generates new config
/remove				remove everything`

	h.sendMessage(chatID, str)
}
