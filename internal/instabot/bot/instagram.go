package bot

import (
	"InstaFollower/internal/pkg/db"

	"github.com/ahmdrz/goinsta/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func subscribeUser(db *db.Database, update tgbotapi.Update, user *goinsta.User) error {
	if err := createInstagramUser(db, user); err != nil {
		return err
	}

	insertFollowersFollowing(db, user)

	if _, err := db.Conn.Exec(sqlInsertSubscription,
		update.Message.Chat.ID,
		user.Username,
	); err != nil {
		return err
	}

	return nil
}

func createInstagramUser(db *db.Database, user *goinsta.User) error {
	if _, err := db.Conn.Exec(sqlInsertInstagramUser,
		user.Username,
		user.FollowerCount,
		user.FollowingCount,
	); err != nil {
		return err
	}

	return nil
}

func insertFollowersFollowing(db *db.Database, user *goinsta.User) {
	followers := user.Followers()
	for followers.Next() {
		for _, follower := range followers.Users {
			db.Conn.Exec(sqlInsertFollowingFollowers,
				follower.Username,
				follower.FullName,
				follower.ExternalURL,
				user.Username,
				"followers",
			)
		}
	}

	following := user.Following()
	for following.Next() {
		for _, following := range following.Users {
			db.Conn.Exec(sqlInsertFollowingFollowers,
				following.Username,
				following.FullName,
				following.ExternalURL,
				user.Username,
				"following",
			)
		}
	}
}

func unsubscribeUser(db *db.Database, update tgbotapi.Update, user *goinsta.User) error {
	if err := deleteInstagramUser(db, user); err != nil {
		return err
	}

	if _, err := db.Conn.Exec(sqlDeleteSubscription,
		update.Message.Chat.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteInstagramUser(db *db.Database, user *goinsta.User) error {
	if _, err := db.Conn.Exec(sqlDeleteInstagramUser,
		user.Username,
	); err != nil {
		return err
	}

	return nil
}
