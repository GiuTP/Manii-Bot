package telegram

import (
	"gfinancer/internal/services"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	api          *tgbotapi.BotAPI
	service      *services.BotService
	allowedUsers map[int64]bool
}

func NewTelegramBot(token string, service *services.BotService, allowedUsers map[int64]bool) (*TelegramBot, error) {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Bot conectado com sucesso na conta: %s", botApi.Self.UserName)
	return &TelegramBot{
		api:          botApi,
		service:      service,
		allowedUsers: allowedUsers,
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

		userID := update.Message.From.ID
		if !b.allowedUsers[userID] {
			log.Printf("ACESSO BLOQUEADO: usuário %s (ID: %d) tentou interagir com o bot", update.Message.From.UserName, userID)
			continue
		}

		textResp, filePath := b.service.HandleMessage(update.Message.Text)

		if filePath != "" {
			msg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
			msg.Caption = textResp

			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Falha ao enviara aquivo para o Telegram: %v", err)
			}

			err := os.Remove(filePath)
			if err != nil {
				log.Printf("Aviso: falha ao apagar arquivo temporário %s: %v", filePath, err)
			}

			continue
		}

		if textResp != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textResp)

			if _, err := b.api.Send(msg); err != nil {
				log.Printf("Falha ao enviar resposta para o Telegram: %v", err)
			}
		}
	}
}
