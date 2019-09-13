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

		// i.send(update.Message.Chat.ID, answers["no_unfollowers"])
	case stateSubscribe:

		user, err := i.instagram.GetUserByUsername(update.Message.Text)
		if err != nil {
			i.send(update.Message.Chat.ID, answers["no_instagram_user"])
			return
		}

		if err := subscribeUser(database, update, user); err != nil {
			i.send(update.Message.Chat.ID, answers["error_user_already_subscribed"])
			createOrUpdateUser(database, update, stateZero)
			return
		}

		createOrUpdateUser(database, update, stateZero)
		i.send(update.Message.Chat.ID, answers["successful_subscription"])

	case stateUnsubscribe:
		defer createOrUpdateUser(database, update, stateZero)

		user, err := getUserForUnsubscription(database, update)
		if err != nil {
			i.send(update.Message.Chat.ID, answers["error_unsubscribe"])
		}

		if update.Message.Text == "Yes" {
			if err := unsubscribeUser(database, update, user); err != nil {
				i.send(update.Message.Chat.ID, answers["error_unsubscribe"])
				return
			}
			
			i.send(update.Message.Chat.ID, answers["successful_unsubscription"])
		} else {
			i.send(update.Message.Chat.ID, answers["unsubscribe_confirmation"])
		}

	default:
		i.send(update.Message.Chat.ID, "кек такого не должно быть")
	}
}
