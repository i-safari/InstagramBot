package bot

import (
	"fmt"
	"log"

	"github.com/Unanoc/InstaFollower/config"
	"github.com/Unanoc/InstaFollower/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	// BASEURL is url of Instagram web site
	BASEURL = "https://www.instagram.com/"
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

	i.Logger.Log("[INFO] ", "Bot is running...")

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		go i.MessageHandler(update)
	}
}

// Send sends message to user
func (i *InstaBot) Send(userID int64, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	i.BotAPI.Send(msg)
}

// Reply replies on message of user
func (i *InstaBot) Reply(userID int64, msgID int, text string) {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyToMessageID = msgID
	i.BotAPI.Send(msg)
}

// MessageHandler ...
func (i *InstaBot) MessageHandler(update tgbotapi.Update) {
	i.Logger.Log("[MESSAGE] ", update.Message.From.UserName, update.Message.Text)

	switch update.Message.Text {
	case "/start":
		i.StartHandler(update)
	default:
		i.MainHandler(update)
	}
}

// StartHandler ...
func (i *InstaBot) StartHandler(update tgbotapi.Update) {
	i.Reply(update.Message.Chat.ID, update.Message.MessageID, STARTANSWER)
}

// MainHandler ...
func (i *InstaBot) MainHandler(update tgbotapi.Update) {
	unfollowers, err := i.InstAccount.GetUnfollowedUsers(update.Message.Text)
	if err != nil {
		i.Send(update.Message.Chat.ID, WRONGUSERNAME)
		return
	}

	length := len(*unfollowers)
	switch {
	case length == 0:
		i.Send(update.Message.Chat.ID, NOUNFOLLOWERS)
		return
	default:
		i.Send(update.Message.Chat.ID, WAITMSG)
		var resultMsg, userInfo string
		var counter int

		for number, user := range *unfollowers {
			userInfo = fmt.Sprintf(
				"%d: %s\nNickname: %s\nFullName: %s\n\n",
				number+1,
				BASEURL+user.Username,
				user.Username,
				user.FullName)

			if counter < 10 {
				resultMsg += userInfo
				counter++
			} else {
				i.Send(update.Message.Chat.ID, resultMsg)
				resultMsg = userInfo
				counter = 1
			}

			if number == len(*unfollowers)-1 {
				i.Send(update.Message.Chat.ID, resultMsg)
			}
		}
		return
	}
}
