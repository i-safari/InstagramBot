package bot

import (
	"InstaFollower/internal/instabot/utils"
	"InstaFollower/internal/pkg/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// InstaBot ...
type InstaBot struct {
	bot      *tgbotapi.BotAPI
	database *db.Database
}

// CreateBot ...
func CreateBot(cfg *utils.Config, database *db.Database) (*InstaBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	if err = getAnswers(cfg); err != nil {
		return nil, err
	}

	return &InstaBot{
		bot:      bot,
		database: database,
	}, nil
}

// Run ...
func (i *InstaBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := i.bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go i.manager(update)
	}
}

// Manager is a router for message handlers
func (i *InstaBot) manager(update tgbotapi.Update) {
	switch update.Message.Text {
	case "/start":
		i.commonHandler(i.database, update, "start", stateZero)
	case "/help":
		i.commonHandler(i.database, update, "help", stateZero)
	case "/listunfollowers":
		i.commonHandler(i.database, update, "list_unfollowers", stateListUnfollowers)
	case "/subscribe":
		i.commonHandler(i.database, update, "subscribe", stateSubscribe)
	case "/unsubscribe":
		i.commonHandler(i.database, update, "unsubscribe", stateUnsubscribe)
	case "/cancel":
		i.cancelHandler(i.database, update)
	default:
		i.statesHandler(i.database, update)
	}
}

// Send sends message to user
func (i *InstaBot) send(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	i.bot.Send(msg)
}
