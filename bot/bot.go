package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/InstaSherlock/config"
	"github.com/InstaSherlock/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	startMsg = `
	Привет! Пришли мне свой username и я расскажу, кто от тебя отписался.
	`
	errorMsg = `
	Произошла какая-то ошибка. Проверь корректность username.
	`
	errorFile = `
	Не удалось сформировать ответ.
	`
)

// InstaBot ...
type InstaBot struct {
	BotAPI      *tgbotapi.BotAPI
	Logger      *Logger
	InstAccount *pkg.Instagram
	Vault       string
}

// CreateBot ...
func CreateBot(cfg *config.Config) (*InstaBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}
	instgram, err := pkg.CreateInstagram(cfg.Username, cfg.Password)
	if err != nil {
		return nil, err
	}

	return &InstaBot{
		BotAPI:      bot,
		Logger:      CreateLogger(cfg.LogPath),
		InstAccount: instgram,
		Vault:       cfg.Vault,
	}, nil
}

// Run ...
func (i *InstaBot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := i.BotAPI.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	i.Logger.LogInfo("Bot is running...")

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		go i.MessageHandler(update)
	}
}

// MessageHandler ...
func (i *InstaBot) MessageHandler(update tgbotapi.Update) {
	i.Logger.Log(update.Message.From.UserName, update.Message.Text)

	switch update.Message.Text {
	case "/start":
		i.StartHandler(update)
	default:
		i.MainHandler(update)
	}
}

// StartHandler ...
func (i *InstaBot) StartHandler(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, startMsg)
	msg.ReplyToMessageID = update.Message.MessageID
	i.BotAPI.Send(msg)
}

// MainHandler ...
func (i *InstaBot) MainHandler(update tgbotapi.Update) {
	userPath := getFilePath(i.Vault, update.Message.Text)
	err := i.InstAccount.GetUnfollowedUsersFile(update.Message.Text, userPath)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
		i.BotAPI.Send(msg)
	} else {
		f, err := os.Open(userPath)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorFile)
			i.BotAPI.Send(msg)
		}
		defer f.Close()
		msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, userPath)
		i.BotAPI.Send(msg)
	}
}

func getFilePath(dir, file string) string {
	if len(dir) == 0 {
		return ""
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}

	if strings.LastIndex(dir, "/") == len(dir)-1 {
		return fmt.Sprintf("%s%s.txt", dir, file)
	}
	return fmt.Sprintf("%s/%s.txt", dir, file)
}
