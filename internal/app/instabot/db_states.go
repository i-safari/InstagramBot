package instabot

import (
	"InstaFollower/internal/pkg/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	stateZERO = iota
	stateLISTUNFOLLOWERS
	stateSUBSCRIBE
	stateUNSUBSCRIBE
)

func getUserState(db *database.Database, update tgbotapi.Update) (state int, err error) {
	err = db.Conn.QueryRow(`
		SELECT "state" 
		FROM 
			users 
		WHERE 
			"id" = $1
	`,
		update.Message.Chat.ID,
	).Scan(&state)

	return
}

func createOrUpdateUser(db *database.Database, update tgbotapi.Update, state int) {
	_, err := getUserState(db, update)

	if err != nil {
		db.Conn.Exec(`
			INSERT INTO "users"
				("id", "username", "firstname", "lastname", "state")
			VALUES 
				($1, $2, $3, $4, $5)
		`,
			update.Message.Chat.ID,
			update.Message.Chat.UserName,
			update.Message.Chat.FirstName,
			update.Message.Chat.LastName,
			state,
		)

		return
	}

	db.Conn.Exec(`
		UPDATE 
			users
		SET 
			"username" = $2,
			"firstname" = $3,
			"lastname" = $4,
			"state" = $5
		WHERE 
			"id" = $1
	`,
		update.Message.Chat.ID,
		update.Message.Chat.UserName,
		update.Message.Chat.FirstName,
		update.Message.Chat.LastName,
		state,
	)
}
