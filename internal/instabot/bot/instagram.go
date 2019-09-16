package bot

import (
	"InstaFollower/internal/pkg/db"
	"fmt"

	"github.com/ahmdrz/goinsta/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func subscribeUser(db *db.Database, update tgbotapi.Update, user *goinsta.User) (err error) {
	if isUserSubscribed(db, update) {
		return fmt.Errorf("subscription already exists")
	}

	if !isInstagramUserExist(db, user) {
		if err = createInstagramUser(db, user); err == nil {
			insertFollowersFollowing(db, user)
		}
	}

	createSubscription(db, update, user)

	return
}

func isUserSubscribed(db *db.Database, update tgbotapi.Update) bool {
	var userID int
	if err := db.Conn.QueryRow(
		sqlSelectSubscription,
		update.Message.Chat.ID,
	).Scan(&userID); err != nil {
		return false
	}

	return true
}

func isInstagramUserExist(db *db.Database, user *goinsta.User) bool {
	var username string
	if err := db.Conn.QueryRow(
		sqlSelectUsername,
		user.Username,
	).Scan(&username); err != nil {
		return false
	}

	return true
}

func createInstagramUser(db *db.Database, user *goinsta.User) (err error) {
	_, err = db.Conn.Exec(sqlInsertInstagramUser,
		user.Username,
		user.FollowerCount,
		user.FollowingCount,
	)

	return
}

func insertFollowersFollowing(db *db.Database, user *goinsta.User) (err error) {
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

	return
}

func createSubscription(db *db.Database, update tgbotapi.Update, user *goinsta.User) (err error) {
	_, err = db.Conn.Exec(sqlInsertSubscription,
		update.Message.Chat.ID,
		user.Username,
	)

	return
}

func unsubscribeUser(db *db.Database, update tgbotapi.Update) error {
	if _, err := db.Conn.Exec(sqlDeleteSubscription,
		update.Message.Chat.ID,
	); err != nil {
		return err
	}

	return nil
}

func deleteSubscription(db *db.Database, update tgbotapi.Update, user *goinsta.User) (err error) {
	_, err = db.Conn.Exec(sqlDeleteSubscription,
		update.Message.Chat.ID,
	)

	return
}
