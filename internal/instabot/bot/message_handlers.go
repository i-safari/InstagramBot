package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (i *InstaBot) commonHandler(
	update tgbotapi.Update,
	answerKey string,
	state int,
) {
	i.send(update.Message.Chat.ID, i.answers[answerKey])
	createOrUpdateUser(i.database, update, state)
}

func (i *InstaBot) cancelHandler(update tgbotapi.Update) {
	userState, err := getUserState(i.database, update)
	if err != nil {
		i.send(update.Message.Chat.ID, i.answers["user_does_not_exist"])
		return
	}

	if userState != stateZERO {
		i.send(update.Message.Chat.ID, i.answers["success_cancel"])
		createOrUpdateUser(i.database, update, stateZERO)
	} else {
		i.send(update.Message.Chat.ID, i.answers["error_cancel"])
	}
}

func (i *InstaBot) statesHandler(update tgbotapi.Update) {
	userState, err := getUserState(i.database, update)
	if err != nil {
		i.send(update.Message.Chat.ID, i.answers["user_does_not_exist"])
		return
	}

	switch userState {

	case stateZERO:

		i.send(update.Message.Chat.ID, i.answers["default"])

	case stateSUBSCRIBE:
		defer createOrUpdateUser(i.database, update, stateZERO)

		user, err := i.instagram.GetUserByUsername(update.Message.Text)
		if err != nil {
			i.send(update.Message.Chat.ID, i.answers["no_instagram_user"])
			return
		}

		if err := subscribeUser(i.database, update, user); err != nil {
			i.send(update.Message.Chat.ID, i.answers["error_user_already_subscribed"])
			return
		}

		i.send(update.Message.Chat.ID, i.answers["successful_subscription"])

	case stateUNSUBSCRIBE:
		// defer createOrUpdateUser(i.database, update, stateZERO)

		// user, err := getUserForUnsubscription(database, update)
		// if err != nil {
		// 	i.send(update.Message.Chat.ID, answers["error_unsubscribe"])
		// }

		// if update.Message.Text == "Yes" {
		// 	if err := unsubscribeUser(database, update, user); err != nil {
		// 		i.send(update.Message.Chat.ID, answers["error_unsubscribe"])
		// 		return
		// 	}

		// 	i.send(update.Message.Chat.ID, answers["successful_unsubscription"])
		// } else {
		// 	i.send(update.Message.Chat.ID, answers["unsubscribe_confirmation"])
		// }

	case stateLISTUNFOLLOWERS:

		// createOrUpdateUser(i.database, update, stateZERO)
		// i.send(update.Message.Chat.ID, answers["no_unfollowers"])

	default:
		i.send(update.Message.Chat.ID, "кек такого не должно быть")
	}
}
