package bot

import "InstaFollower/internal/pkg/db"

const (
	stateZero = iota
	stateListUnfollowers
	stateSubscribe
	stateUnsubscribe
)

func getUserState(db *db.Database, userID int64) (state int) {
	db.Conn.QueryRow(sqlSelectUserStateByUserID, userID).Scan(&state)
	return // if user's state is found in database, it will returns 0 by default
}

func createOrUpdateUserState(db *db.Database, userID int64, state int) {
	if _, err := db.Conn.Exec(sqlInsertUserState, userID, state); err != nil {
		db.Conn.Exec(sqlUpdateUserState, userID, state)
	}
}
