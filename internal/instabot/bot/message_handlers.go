package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (i *InstaBot) commonHandler(update tgbotapi.Update, commandName string) {
	i.send(update.Message.Chat.ID, answers[commandName])
}

// func (i *InstaBot) subscribeHandler(update tgbotapi.Update) {

// }

// func (i *InstaBot) unsubscribeHandler(update tgbotapi.Update) {

// }


// func (i *InstaBot) getListOfUnfollowers()