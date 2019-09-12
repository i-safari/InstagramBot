package bot

import (
	"InstaFollower/internal/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (i *InstaBot) commonHandler(
	database *db.Database,
	update tgbotapi.Update,
	answerKey string,
	state int,
) {
	i.send(update.Message.Chat.ID, answers[answerKey])

	createOrUpdateUserState(database, update.Message.Chat.ID, state)
}

func (i *InstaBot) cancelHandler(
	database *db.Database,
	update tgbotapi.Update,
) {
	userState := getUserState(database, update.Message.Chat.ID)

	if userState != stateZero {
		i.send(update.Message.Chat.ID, answers["success_cancel"])
		createOrUpdateUserState(database, update.Message.Chat.ID, stateZero)
	} else {
		i.send(update.Message.Chat.ID, answers["error_cancel"])
	}
}

func (i *InstaBot) statesHandler(
	database *db.Database,
	update tgbotapi.Update,
) {
	userState := getUserState(i.database, update.Message.Chat.ID)

	switch userState {
	case stateZero:
		i.send(update.Message.Chat.ID, answers["default"])
	case stateListUnfollowers:
		i.send(update.Message.Chat.ID, "(отдаю список отписчиков)")
		createOrUpdateUserState(database, update.Message.Chat.ID, stateZero)
	case stateSubscribe:
		i.send(update.Message.Chat.ID, "(подписываю юзера)")
		createOrUpdateUserState(database, update.Message.Chat.ID, stateZero)
	case stateUnsubscribe:
		if update.Message.Text == "Yes" {
			i.send(update.Message.Chat.ID, "(отписываю юзера)")
			createOrUpdateUserState(database, update.Message.Chat.ID, stateZero)
		} else {
			i.send(update.Message.Chat.ID, answers["unsubscribeconfirmation"])
			createOrUpdateUserState(database, update.Message.Chat.ID, stateZero)
		}
	default:
		i.send(update.Message.Chat.ID, "кек такого не должно быть")
	}
}
