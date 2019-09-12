package bot

// Table `state`
const (
	sqlSelectUserStateByUserID = `
		SELECT "state" 
		FROM states 
		WHERE "user_id" = $1
	`
	sqlInsertUserState = `
		INSERT INTO states
		("user_id", "state")
		VALUES ($1, $2)
	`
	sqlUpdateUserState = `
		UPDATE states
		SET "state" = $2
		WHERE "user_id" = $1
	`
)
