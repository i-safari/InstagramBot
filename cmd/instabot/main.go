package main

import (
	"InstaFollower/internal/instabot/bot"
)

func main() {
	bot, err := bot.CreateBot()
	if err != nil {
		panic(err)
	}

	bot.Run()
}
