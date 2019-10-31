package instabot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (i *InstaBot) commonHandler(
	update tgbotapi.Update,
	answerKey string,
	state int,
) {
	i.bot.Send(update.Message.Chat.ID, i.answers[answerKey])
	createOrUpdateUser(i.db, update, state)
}

func (i *InstaBot) cancelHandler(update tgbotapi.Update) {
	userState, err := getUserState(i.db, update)
	if err != nil {
		i.bot.Send(update.Message.Chat.ID, i.answers["USER_does_not_exist_error"])
		return
	}

	if userState != stateZERO {
		i.bot.Send(update.Message.Chat.ID, i.answers["CANCEL_ok"])
		createOrUpdateUser(i.db, update, stateZERO)
	} else {
		i.bot.Send(update.Message.Chat.ID, i.answers["CANCEL_error"])
	}
}

func (i *InstaBot) statesHandler(update tgbotapi.Update) {
	userState, err := getUserState(i.db, update)
	if err != nil {
		i.bot.Send(update.Message.Chat.ID, i.answers["USER_does_not_exist_error"])
		return
	}

	switch userState {

	case stateZERO:
		i.bot.Send(update.Message.Chat.ID, i.answers["DEFAULT"])
	case stateSUBSCRIBE:
		defer createOrUpdateUser(i.db, update, stateZERO)
		subscribeUser(i, update)
	case stateUNSUBSCRIBE:
		defer createOrUpdateUser(i.db, update, stateZERO)
		unsubscribeUser(i, update)
	case stateLISTUNFOLLOWERS:
		defer createOrUpdateUser(i.db, update, stateZERO)
		getListUnfollowers(i, update)
	}
}
