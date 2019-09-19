package bot

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx"
)

func subscribeUser(i *InstaBot, update tgbotapi.Update) {
	user, err := i.instagram.GetUserByUsername(update.Message.Text)
	if err != nil {
		i.Send(update.Message.Chat.ID, i.answers["no_instagram_user"])
		return
	}

	if user.IsPrivate {
		i.Send(update.Message.Chat.ID, i.answers["instagram_user_is_private"])
		return
	}

	if isUserSubscribed(i.database, update) {
		i.Send(update.Message.Chat.ID, i.answers["error_user_already_subscribed"])
		return
	}

	if !isInstagramUserExist(i.database, user) {
		if err = createInstagramUser(i.database, user); err == nil {
			createFollowersFollowing(i.database, user)
		}
	}

	createSubscription(i.database, update, user)

	i.Send(update.Message.Chat.ID, i.answers["successful_subscription"])
}

func unsubscribeUser(i *InstaBot, update tgbotapi.Update) {
	if !isUserSubscribed(i.database, update) {
		i.Send(update.Message.Chat.ID, i.answers["error_unsubscribe"])
		return
	}

	if update.Message.Text == "Yes" {
		if err := deleteSubscription(i.database, update); err != nil {
			i.Send(update.Message.Chat.ID, i.answers["error_unsubscribe"])
			return
		}

		i.Send(update.Message.Chat.ID, i.answers["successful_unsubscription"])
	} else {
		i.Send(update.Message.Chat.ID, i.answers["unsubscribe_confirmation_fail"])
	}
}

func getListUnfollowers(i *InstaBot, update tgbotapi.Update) {
	user, err := i.instagram.GetUserByUsername(update.Message.Text)
	if err != nil {
		i.Send(update.Message.Chat.ID, i.answers["no_instagram_user"])
		return
	}

	if user.IsPrivate {
		i.Send(update.Message.Chat.ID, i.answers["instagram_user_is_private"])
		return
	}

	if !isInstagramUserExist(i.database, user) {
		if err = createInstagramUser(i.database, user); err == nil {
			createFollowersFollowing(i.database, user)
		}
	}

	unfollowers, err := selectUnfollowers(i.database, user.Username)
	if err != nil {
		i.Send(update.Message.Chat.ID, i.answers["failed_unfollowers_list"])
		return
	}

	pathToCSV := fmt.Sprintf("%s/%d.csv", i.vaultPath, update.Message.Chat.ID)

	err = makeCSV(unfollowers, update, pathToCSV)
	if err != nil {
		i.Send(update.Message.Chat.ID, i.answers["failed_unfollowers_list"])
		return
	}

	i.SendDocument(update.Message.Chat.ID, pathToCSV)
	_ = os.Remove(pathToCSV)
}

func makeCSV(unfollowers *pgx.Rows, update tgbotapi.Update, pathCSV string) error {
	f, err := os.Create(pathCSV)
	if err != nil {
		return err
	}
	defer f.Close()

	var title = "USERNAME,FULLNAME,LINK\n"

	if _, err = f.WriteString(title); err != nil {
		return err
	}

	for unfollowers.Next() {
		var username, fullname, url string

		_ = unfollowers.Scan(&username, &fullname, &url)

		csvString := fmt.Sprintf("%s,%s,%s\n", username, fullname, url)

		if _, err := f.WriteString(csvString); err != nil {
			return err
		}
	}

	return nil
}
