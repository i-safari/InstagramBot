package main

import bot "InstaFollower/internal/instabot"

func main() {
	bot, err := bot.CreateBot()
	if err != nil {
		panic(err)
	}

	bot.Run()
}
