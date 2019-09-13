package bot

import (
	"InstaFollower/internal/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	stateZero = iota
	stateListUnfollowers
	stateSubscribe
	stateUnsubscribe
)

func getUserState(db *db.Database, update tgbotapi.Update) (state int) {
	db.Conn.QueryRow(sqlSelectUserStateByUserID, update.Message.Chat.ID).Scan(&state)
	return // if user's state is found in database, it will returns 0 by default
}

func createOrUpdateUser(db *db.Database, update tgbotapi.Update, state int) {
	if _, err := db.Conn.Exec(sqlInsertUser,
		update.Message.Chat.ID,
		update.Message.Chat.UserName,
		update.Message.Chat.FirstName,
		update.Message.Chat.LastName,
		state,
	); err != nil {
		db.Conn.Exec(sqlUpdateUser,
			update.Message.Chat.ID,
			update.Message.Chat.UserName,
			update.Message.Chat.FirstName,
			update.Message.Chat.LastName,
			state,
		)
	}
}
