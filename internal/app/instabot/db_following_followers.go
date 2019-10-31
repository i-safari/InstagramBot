package instabot

import (
	"InstaFollower/internal/pkg/database"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/jackc/pgx"
)

func createFollowersFollowing(db *database.Database, user *goinsta.User) (err error) {
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

func selectUnfollowers(db *database.Database, username string) (unfollowers *pgx.Rows, err error) {
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
