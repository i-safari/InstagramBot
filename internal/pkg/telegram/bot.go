package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CreateBot create an instance of Bot
func CreateBot(telegramToken string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	return &Bot{api}, nil
}

// Bot is a custom wrapper for Telegram Bot
type Bot struct {
	API *tgbotapi.BotAPI
}

// GetUpdatesChanel returns UpdatesChannel for receiving new updates
func (i *Bot) GetUpdatesChanel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := i.API.GetUpdatesChan(u)
	return updates, err
}

// Send sends message to user
func (i *Bot) Send(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	i.API.Send(msg)
}

// SendDocument sends a document to user
func (i *Bot) SendDocument(userID int64, file interface{}) {
	msg := tgbotapi.NewDocumentUpload(userID, file)
	i.API.Send(msg)
}
