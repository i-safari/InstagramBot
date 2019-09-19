package bot

import (
	"InstaFollower/internal/pkg/db"

	"github.com/ahmdrz/goinsta/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx"
)

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

func createFollowersFollowing(db *db.Database, user *goinsta.User) (err error) {
	var baseURL = "https://instagram.com/"

	followers := user.Followers()
	for followers.Next() {
		for _, follower := range followers.Users {
			db.Conn.Exec(sqlInsertFollowingFollowers,
				follower.Username,
				follower.FullName,
				baseURL+follower.Username,
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
				baseURL+following.Username,
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

func deleteSubscription(db *db.Database, update tgbotapi.Update) (err error) {
	_, err = db.Conn.Exec(sqlDeleteSubscription,
		update.Message.Chat.ID,
	)

	return
}

func selectUnfollowers(db *db.Database, username string) (unfollowers *pgx.Rows, err error) {
	unfollowers, err = db.Conn.Query(sqlSelectUnfollowers, username)
	return
}
