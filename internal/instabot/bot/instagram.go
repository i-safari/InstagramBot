package bot

import (
	"InstaFollower/internal/pkg/db"
	"fmt"

	"github.com/ahmdrz/goinsta/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func subscribeUser(db *db.Database, update tgbotapi.Update, user *goinsta.User) (err error) {
	if isUserSubscribed(db, update, user) { // почеум не работает???
		return fmt.Errorf("subscription already exists")
	}

	if !isInstagramUserExist(db, user) {
		createInstagramUser(db, user)
		insertFollowersFollowing(db, user)
	}

	createSubscription(db, update, user)

	return
}

func isUserSubscribed(db *db.Database, update tgbotapi.Update, user *goinsta.User) bool {
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
		sqlSelectInsgramUser,
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

// func unsubscribeUser(db *db.Database, update tgbotapi.Update, user *goinsta.User) error {
// 	if err := deleteInstagramUser(db, user); err != nil {
// 		return err
// 	}

// 	if _, err := db.Conn.Exec(sqlDeleteSubscription,
// 		update.Message.Chat.ID,
// 	); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func deleteInstagramUser(db *db.Database, user *goinsta.User) error {
// 	if _, err := db.Conn.Exec(sqlDeleteInstagramUser,
// 		user.Username,
// 	); err != nil {
// 		return err
// 	}

// 	return nil
// }
