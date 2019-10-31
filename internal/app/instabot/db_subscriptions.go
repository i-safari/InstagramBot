package instabot

import (
	"InstaFollower/internal/pkg/database"

	"github.com/ahmdrz/goinsta/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func createSubscription(db *database.Database, update tgbotapi.Update, user *goinsta.User) (err error) {
	_, err = db.Conn.Exec(`
		INSERT INTO "subscriptions"
			("user_id", "insta_username")
		VALUES 
			($1, $2)
	`,
		update.Message.Chat.ID,
		user.Username,
	)

	return
}

func deleteSubscription(db *database.Database, update tgbotapi.Update) (err error) {
	_, err = db.Conn.Exec(`
		DELETE 
		FROM 
			"subscriptions"
		WHERE "user_id" = $1
	`,
		update.Message.Chat.ID,
	)

	return
}

func isUserSubscribed(db *database.Database, update tgbotapi.Update) bool {
	var userID int
	if err := db.Conn.QueryRow(`
		SELECT 
			"user_id"
		FROM 
			"subscriptions"
		WHERE
			"user_id" = $1
	`,
		update.Message.Chat.ID,
	).Scan(&userID); err != nil {
		return false
	}

	return true
}
