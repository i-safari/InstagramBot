package bot

// Table `state`
const (
	sqlSelectUserStateByUserID = `
		SELECT "state" 
		FROM users 
		WHERE "id" = $1
	`
	sqlInsertUser = `
		INSERT INTO users
		("id", "username", "firstname", "lastname", "state")
		VALUES ($1, $2, $3, $4, $5)
	`
	sqlUpdateUser = `
		UPDATE users
		SET 
			"username" = $2,
			"firstname" = $3,
			"lastname" = $4,
			"state" = $5
		WHERE "id" = $1
	`
)
