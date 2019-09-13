package bot

// Table "users"
const (
	sqlSelectUserStateByUserID = `
		SELECT "state" 
		FROM 
			users 
		WHERE 
			"id" = $1
	`
	sqlInsertUser = `
		INSERT INTO "users"
			("id", "username", "firstname", "lastname", "state")
		VALUES 
			($1, $2, $3, $4, $5)
	`
	sqlUpdateUser = `
		UPDATE 
			users
		SET 
			"username" = $2,
			"firstname" = $3,
			"lastname" = $4,
			"state" = $5
		WHERE 
			"id" = $1
	`
)

// Table "subscriptions"
const (
	sqlInsertSubscription = `
		INSERT INTO "subscriptions"
			("user_id", "insta_username")
		VALUES 
			($1, $2)
	`
	sqlDeleteSubscription = `
		DELETE 
		FROM 
			"subscriptions"
		WHERE "user_id" = ($1)
	`
)

// Table "insta_users"
const (
	sqlInsertInstagramUser = `
		INSERT INTO "insta_users"
			("username", "followers", "following")
		VALUES 
			($1, $2, $3)
	`
	sqlDeleteInstagramUser = `
		DELETE 
		FROM 
			"insta_users"
		WHERE
			"username" = $1
	`
)

// Table "following_followers"
const (
	sqlInsertFollowingFollowers = `
		INSERT INTO "following_followers"
			("username", "fullname", "URL", "refer_username", "group_type")
		VALUES 
			($1, $2, $3, $4, $5)
	`
)
