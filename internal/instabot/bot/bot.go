package bot

import (
	"InstaFollower/internal/instabot/instagram"
	"InstaFollower/internal/pkg/db"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// InstaBot ...
type InstaBot struct {
	bot       *tgbotapi.BotAPI
	database  *db.Database
	instagram *instagram.Instagram
	answers   map[string]string
	vaultPath string
}

// CreateBot ...
func CreateBot() (*InstaBot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return nil, err
	}

	answers, err := getAnswers(os.Getenv("PATH_TO_DIALOGS"))
	if err != nil {
		return nil, err
	}

	instagram, err := instagram.CreateInstagram(os.Getenv("LOGIN"), os.Getenv("PASSWORD"))
	if err != nil {
		return nil, err
	}

	database, err := db.CreateConnection(os.Getenv("POSTGRES_URI"))
	if err != nil {
		return nil, err
	}

	path := os.Getenv("VAULT_PATH")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	return &InstaBot{
		bot:       bot,
		database:  database,
		instagram: instagram,
		answers:   answers,
		vaultPath: path,
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

		go i.Manager(update)
	}
}

// Manager is a router for message handlers
func (i *InstaBot) Manager(update tgbotapi.Update) {
	switch update.Message.Text {
	case "/start":
		i.commonHandler(update, "start", stateZERO)
	case "/help":
		i.commonHandler(update, "help", stateZERO)
	case "/listunfollowers":
		i.commonHandler(update, "list_unfollowers", stateLISTUNFOLLOWERS)
	case "/subscribe":
		i.commonHandler(update, "subscribe", stateSUBSCRIBE)
	case "/unsubscribe":
		i.commonHandler(update, "unsubscribe", stateUNSUBSCRIBE)
	case "/cancel":
		i.cancelHandler(update)
	default:
		i.statesHandler(update)
	}
}

// Send sends message to user
func (i *InstaBot) Send(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	i.bot.Send(msg)
}

// SendDocument sends a document to user
func (i *InstaBot) SendDocument(userID int64, file interface{}) {
	msg := tgbotapi.NewDocumentUpload(userID, file)
	i.bot.Send(msg)
}
