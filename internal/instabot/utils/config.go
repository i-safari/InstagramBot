package utils

import "os"

// Config is config struct for instabot.
type Config struct {
	Login         string
	Password      string
	TelegramToken string
	PostgresURI   string
	PathToDialogs string
}

// CreateConfig returns initialized config.
func CreateConfig() *Config {
	var config Config

	config.Login = os.Getenv("LOGIN")
	config.Password = os.Getenv("PASSWORD")
	config.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	config.PostgresURI = os.Getenv("POSTGRES_URI")
	config.PathToDialogs = os.Getenv("PATH_TO_DIALOGS")

	return &config
}
