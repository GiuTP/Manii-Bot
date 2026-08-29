package telegram

import (
	"gfinancer/internal/services"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	api     *tgbotapi.BotAPI
	service *services.BotService
}

func NewTelegramBot(token string, service *services.BotService) (*TelegramBot, error) {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Bot conectado com sucesso na conta: %s", botApi.Self.UserName)
	return &TelegramBot{
		api:     botApi,
		service: service,
	}, nil
}

func (b *TelegramBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		resp := b.service.HandleMessage(update.Message.Text)

		if resp != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
			b.api.Send(msg)
		}
	}
}
