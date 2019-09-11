package main

import (
	"InstaFollower/internal/instabot/bot"
	"InstaFollower/internal/instabot/db"
	"InstaFollower/internal/instabot/utils"
)

func main() {
	config := utils.CreateConfig()

	db, err := db.CreateConnection(config.PostgresURI)
	if err != nil {
		panic(err)
	}
	defer db.Disconnect()

	bot, err := bot.CreateBot(config, db)
	if err != nil {
		panic(err)
	}

	bot.Run()
}
