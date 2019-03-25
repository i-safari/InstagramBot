package main

import (
	"flag"
	"log"

	"github.com/Unanoc/InstaFollower/bot"
	"github.com/Unanoc/InstaFollower/config"
)

var (
	configPath = flag.String("cfg", "", "Path to configuration JSON file")
)

func main() {
	flag.Parse()

	// getting config
	cfg, err := config.CreateConfig(*configPath)
	if err != nil {
		log.Panic(err)
	}

	// bot initializing
	instabot, err := bot.CreateBot(cfg)
	if err != nil {
		log.Panic(err)
	}

	// starting the bot
	instabot.Run()
}
