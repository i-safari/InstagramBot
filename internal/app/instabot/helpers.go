package instabot

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx"
)

func subscribeUser(i *InstaBot, update tgbotapi.Update) {
	user, err := i.instagram.GetUserByUsername(update.Message.Text)
	if err != nil {
		i.bot.Send(update.Message.Chat.ID, i.answers["INSTAGRAM_user_not_found"])
		log.Printf("[ERROR] function subscribeUser(): %s", err)
		return
	}

	if user.IsPrivate {
		i.bot.Send(update.Message.Chat.ID, i.answers["INSTAGRAM_user_is_private"])
		return
	}

	if isUserSubscribed(i.db, update) {
		i.bot.Send(update.Message.Chat.ID, i.answers["SUBSCRIBE_user_already_subscribed_error"])
		return
	}

	if !isInstagramUserExist(i.db, user) {
		if err = createInstagramUser(i.db, user); err == nil {
			createFollowersFollowing(i.db, user)
		}
	}

	createSubscription(i.db, update, user)

	i.bot.Send(update.Message.Chat.ID, i.answers["SUBSCRIPTION_ok"])
}

func unsubscribeUser(i *InstaBot, update tgbotapi.Update) {
	if !isUserSubscribed(i.db, update) {
		i.bot.Send(update.Message.Chat.ID, i.answers["UNSUBSCRIBE_error"])
		return
	}

	if update.Message.Text == "Yes" {
		if err := deleteSubscription(i.db, update); err != nil {
			i.bot.Send(update.Message.Chat.ID, i.answers["UNSUBSCRIBE_error"])
			log.Printf("[ERROR] function unsubscribeUser(): %s", err)
			return
		}

		i.bot.Send(update.Message.Chat.ID, i.answers["UNSUBSCRIBE_ok"])
	} else {
		i.bot.Send(update.Message.Chat.ID, i.answers["UNSUBSCRIBE_confirmation_error"])
	}
}

func getListUnfollowers(i *InstaBot, update tgbotapi.Update) {
	user, err := i.instagram.GetUserByUsername(update.Message.Text)
	if err != nil {
		i.bot.Send(update.Message.Chat.ID, i.answers["INSTAGRAM_user_not_found"])
		log.Printf("[ERROR] function getListUnfollowers(): %s", err)
		return
	}

	if user.IsPrivate {
		i.bot.Send(update.Message.Chat.ID, i.answers["INSTAGRAM_user_is_private"])
		log.Printf("[ERROR] function getListUnfollowers(): %s", err)
		return
	}

	if !isInstagramUserExist(i.db, user) {
		if err = createInstagramUser(i.db, user); err == nil {
			createFollowersFollowing(i.db, user)
		}
	}

	unfollowers, err := selectUnfollowers(i.db, user.Username)
	if err != nil {
		i.bot.Send(update.Message.Chat.ID, i.answers["UNFOLLOWERS_LIST_error"])
		log.Printf("[ERROR] function getListUnfollowers(): %s", err)
		return
	}

	pathToCSV := fmt.Sprintf("%s/%d.csv", i.vaultPath, update.Message.Chat.ID)

	err = makeCSV(unfollowers, update, pathToCSV)
	if err != nil {
		i.bot.Send(update.Message.Chat.ID, i.answers["UNFOLLOWERS_LIST_error"])
		log.Printf("[ERROR] function getListUnfollowers(): %s", err)
		return
	}

	i.bot.SendDocument(update.Message.Chat.ID, pathToCSV)
	_ = os.Remove(pathToCSV)
}

func makeCSV(unfollowers *pgx.Rows, update tgbotapi.Update, pathCSV string) error {
	f, err := os.Create(pathCSV)
	if err != nil {
		log.Printf("[ERROR] function makeCSV(): %s", err)
		return err
	}
	defer f.Close()

	var title = "USERNAME,FULLNAME,LINK\n"

	if _, err = f.WriteString(title); err != nil {
		log.Printf("[ERROR] function makeCSV(): %s", err)
		return err
	}

	for unfollowers.Next() {
		var username, fullname, url string

		_ = unfollowers.Scan(&username, &fullname, &url)

		csvString := fmt.Sprintf("%s,%s,%s\n", username, fullname, url)

		if _, err := f.WriteString(csvString); err != nil {
			log.Printf("[ERROR] function makeCSV(): %s", err)
			return err
		}
	}

	return nil
}
