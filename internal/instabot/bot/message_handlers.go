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

	createOrUpdateUser(database, update, state)
}

func (i *InstaBot) cancelHandler(
	database *db.Database,
	update tgbotapi.Update,
) {
	userState := getUserState(database, update)

	if userState != stateZero {
		i.send(update.Message.Chat.ID, answers["success_cancel"])
		createOrUpdateUser(database, update, stateZero)
	} else {
		i.send(update.Message.Chat.ID, answers["error_cancel"])
	}
}

func (i *InstaBot) statesHandler(
	database *db.Database,
	update tgbotapi.Update,
) {
	userState := getUserState(i.database, update)

	switch userState {
	case stateZero:
		i.send(update.Message.Chat.ID, answers["default"])
	case stateListUnfollowers:
		// unfollowers, err := instagram.GetUnfollowers(update.Message.Text)
		// if err != nil {
		// 	i.send(update.Message.Chat.ID, answers["error_private_account"])
		// 	createOrUpdateUser(database, update.Message.Chat.ID, stateZero)
		// 	return
		// }

		// if len(unfollowers) > 0 {
		// 	// make csv
		// 	// send csv
		// 	createOrUpdateUser(database, update.Message.Chat.ID, stateZero)
		// }

		i.send(update.Message.Chat.ID, answers["no_unfollowers"])
	case stateSubscribe:
		i.send(update.Message.Chat.ID, "(подписываю юзера)")
		createOrUpdateUser(database, update, stateZero)
	case stateUnsubscribe:
		if update.Message.Text == "Yes" {
			i.send(update.Message.Chat.ID, "(отписываю юзера)")
			createOrUpdateUser(database, update, stateZero)
		} else {
			i.send(update.Message.Chat.ID, answers["unsubscribe_confirmation"])
		}
	default:
		i.send(update.Message.Chat.ID, "кек такого не должно быть")
	}
}
