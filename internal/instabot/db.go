package bot

import (
	"InstaFollower/internal/pkg/db"

	"github.com/ahmdrz/goinsta/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx"
)

func isUserSubscribed(db *db.Database, update tgbotapi.Update) bool {
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

func isInstagramUserExist(db *db.Database, user *goinsta.User) bool {
	var username string
	if err := db.Conn.QueryRow(`
		SELECT
			"username"
		FROM
			"insta_users"
		WHERE
			"username" = $1
	`,
		user.Username,
	).Scan(&username); err != nil {
		return false
	}

	return true
}

func createInstagramUser(db *db.Database, user *goinsta.User) (err error) {
	_, err = db.Conn.Exec(`
		INSERT INTO "insta_users"
			("username", "followers", "following")
		VALUES 
			($1, $2, $3)
	`,
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
			db.Conn.Exec(`
				INSERT INTO "following_followers"
					("username", "fullname", "URL", "refer_username", "group_type")
				VALUES 
					($1, $2, $3, $4, $5)
			`,
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
			db.Conn.Exec(`
				INSERT INTO "following_followers"
					("username", "fullname", "URL", "refer_username", "group_type")
				VALUES 
					($1, $2, $3, $4, $5)
			`,
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

func deleteSubscription(db *db.Database, update tgbotapi.Update) (err error) {
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

func selectUnfollowers(db *db.Database, username string) (unfollowers *pgx.Rows, err error) {
	unfollowers, err = db.Conn.Query(`
		SELECT
			"username",
			"fullname",
			"URL"
		FROM
			"following_followers"
		WHERE
			"group_type" = 'following' AND
			"refer_username" = $1 AND
			"username" NOT IN (
				SELECT
					"username"
				FROM
					"following_followers"
				WHERE
					"group_type" = 'followers' AND
					"refer_username" = $1
			)
		ORDER BY
			"fullname"
	`,
		username,
	)
	return
}
