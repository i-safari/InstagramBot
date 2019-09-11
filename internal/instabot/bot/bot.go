package bot

import (
	"InstaFollower/internal/instabot/db"
	"InstaFollower/internal/instabot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// InstaBot ...
type InstaBot struct {
	bot  *tgbotapi.BotAPI
	conn *db.Database
}

// CreateBot ...
func CreateBot(cfg *utils.Config, conn *db.Database) (*InstaBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	if err = getAnswers(cfg); err != nil {
		return nil, err
	}

	return &InstaBot{
		bot:  bot,
		conn: conn,
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
		i.commonHandler(update, "start")
	case "/help":
		i.commonHandler(update, "help")
	// case "/subscribe":
	// 	i.subscribeHandler(update)
	// case "/unsubscribe":
	// 	i.unsubscribeHandler(update)
	default:
		i.commonHandler(update, "default")
	}
}

// Send sends message to user
func (i *InstaBot) send(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	i.bot.Send(msg)
}

// Reply replies on message of user
func (i *InstaBot) reply(userID int64, msgID int, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyToMessageID = msgID
	i.bot.Send(msg)
}
