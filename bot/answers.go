package bot

const (
	// STARTANSWER is answer for message "/start"
	STARTANSWER = `
	Привет! Пришли мне свой username в Instagram без "@" и я расскажу, кто не подписан на тебя в ответ.
	Но только, если у тебя не приватный аккаунт!
	`
	// WRONGUSERNAME is answer for incorrect username
	WRONGUSERNAME = `
	Произошла какая-то ошибка. Проверь корректность username.
	`
	// NOUNFOLLOWERS is answer for case when unfollowers list is empty
	NOUNFOLLOWERS = `
	Вау, да у тебя нет таких людей, которые не подписаны на тебя в ответ!
	`
	// WAITMSG is message before sending unfollowers
	WAITMSG = `
	Пожалуйста подожди, это займет какое-то время:)
	`
)
