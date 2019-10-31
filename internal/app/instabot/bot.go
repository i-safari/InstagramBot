package instabot

import (
	"InstaFollower/internal/pkg/database"
	"InstaFollower/internal/pkg/instagram"
	"InstaFollower/internal/pkg/telegram"
	"InstaFollower/pkg/utils"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// InstaBot ...
type InstaBot struct {
	bot       *telegram.Bot
	db        *database.Database
	instagram *instagram.Instagram
	answers   map[string]string
	vaultPath string
}

// CreateBot ...
func CreateBot() (*InstaBot, error) {
	bot, err := telegram.CreateBot(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return nil, err
	}

	db, err := database.CreateConnection(os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, err
	}

	instagram, err := instagram.CreateInstagram(os.Getenv("LOGIN"), os.Getenv("PASSWORD"))
	if err != nil {
		return nil, err
	}

	answers, err := utils.GetMapFromJSON(os.Getenv("PATH_TO_DIALOGS"))
	if err != nil {
		return nil, err
	}

	path := os.Getenv("VAULT_PATH")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	return &InstaBot{
		bot:       bot,
		db:        db,
		instagram: instagram,
		answers:   answers,
		vaultPath: path,
	}, nil
}

// Run ...
func (i *InstaBot) Run() {
	log.Println("[INFO] bot is running...")

	updates, err := i.bot.GetUpdatesChanel()
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go i.Manager(update)

		log.Printf("[UPDATE] user_id=%d username=%s first_name=%s last_name=%s msg=%s",
			update.Message.Chat.ID,
			update.Message.Chat.UserName,
			update.Message.Chat.FirstName,
			update.Message.Chat.LastName,
			update.Message.Text,
		)
	}
}

// Manager is a router for message handlers
func (i *InstaBot) Manager(update tgbotapi.Update) {
	switch update.Message.Text {
	case "/start":
		i.commonHandler(update, update.Message.Text, stateZERO)
	case "/help":
		i.commonHandler(update, update.Message.Text, stateZERO)
	case "/listunfollowers":
		i.commonHandler(update, update.Message.Text, stateLISTUNFOLLOWERS)
	case "/subscribe":
		i.commonHandler(update, update.Message.Text, stateSUBSCRIBE)
	case "/unsubscribe":
		i.commonHandler(update, update.Message.Text, stateUNSUBSCRIBE)
	case "/cancel":
		i.cancelHandler(update)
	default:
		i.statesHandler(update)
	}
}
