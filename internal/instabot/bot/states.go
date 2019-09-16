package bot

import (
	"InstaFollower/internal/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	stateZERO = iota
	stateLISTUNFOLLOWERS
	stateSUBSCRIBE
	stateUNSUBSCRIBE
)

func getUserState(db *db.Database, update tgbotapi.Update) (state int, err error) {
	err = db.Conn.QueryRow(sqlSelectUserStateByUserID, update.Message.Chat.ID).Scan(&state)
	return
}

func createOrUpdateUser(db *db.Database, update tgbotapi.Update, state int) {
	_, err := getUserState(db, update)

	if err != nil {
		db.Conn.Exec(sqlInsertUser,
			update.Message.Chat.ID,
			update.Message.Chat.UserName,
			update.Message.Chat.FirstName,
			update.Message.Chat.LastName,
			state,
		)

		return
	}

	db.Conn.Exec(sqlUpdateUser,
		update.Message.Chat.ID,
		update.Message.Chat.UserName,
		update.Message.Chat.FirstName,
		update.Message.Chat.LastName,
		state,
	)
}
