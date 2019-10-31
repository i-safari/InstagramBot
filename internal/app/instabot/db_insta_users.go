package instabot

import (
	"InstaFollower/internal/pkg/database"

	"github.com/ahmdrz/goinsta/v2"
)

func createInstagramUser(db *database.Database, user *goinsta.User) (err error) {
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

func isInstagramUserExist(db *database.Database, user *goinsta.User) bool {
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
