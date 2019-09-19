package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (i *InstaBot) commonHandler(
	update tgbotapi.Update,
	answerKey string,
	state int,
) {
	i.Send(update.Message.Chat.ID, i.answers[answerKey])
	createOrUpdateUser(i.database, update, state)
}

func (i *InstaBot) cancelHandler(update tgbotapi.Update) {
	userState, err := getUserState(i.database, update)
	if err != nil {
		i.Send(update.Message.Chat.ID, i.answers["user_does_not_exist"])
		return
	}

	if userState != stateZERO {
		i.Send(update.Message.Chat.ID, i.answers["success_cancel"])
		createOrUpdateUser(i.database, update, stateZERO)
	} else {
		i.Send(update.Message.Chat.ID, i.answers["error_cancel"])
	}
}

func (i *InstaBot) statesHandler(update tgbotapi.Update) {
	userState, err := getUserState(i.database, update)
	if err != nil {
		i.Send(update.Message.Chat.ID, i.answers["user_does_not_exist"])
		return
	}

	switch userState {

	case stateZERO:
		i.Send(update.Message.Chat.ID, i.answers["default"])
	case stateSUBSCRIBE:
		defer createOrUpdateUser(i.database, update, stateZERO)
		subscribeUser(i, update)
	case stateUNSUBSCRIBE:
		defer createOrUpdateUser(i.database, update, stateZERO)
		unsubscribeUser(i, update)
	case stateLISTUNFOLLOWERS:
		defer createOrUpdateUser(i.database, update, stateZERO)
		getListUnfollowers(i, update)
	}
}
