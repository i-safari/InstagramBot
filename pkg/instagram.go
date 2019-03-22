package pkg

import (
	"log"
	"os"

	"gopkg.in/ahmdrz/goinsta.v2"
)

// CreateInstagram ...
func CreateInstagram(username, password string) (*Instagram, error) {
	var instagram Instagram
	if err := instagram.SetAccount(username, password); err != nil {
		log.Println(err)
		return nil, err
	}

	return &instagram, nil
}

// Instagram ...
type Instagram struct {
	Account *goinsta.Instagram
}

// SetAccount ...
func (i *Instagram) SetAccount(username, password string) error {
	i.Account = goinsta.New(username, password)
	if err := i.Account.Login(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetUserByUsername ...
func (i *Instagram) GetUserByUsername(username string) (*goinsta.User, error) {
	user, err := i.Account.Profiles.ByName(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

// GetUnfollowedUsers ...
func (i *Instagram) GetUnfollowedUsers(username string) (*[]string, error) {
	user, err := i.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	followersMap := make(map[string]bool, 0)
	followers := user.Followers()
	for followers.Next() {
		for _, user := range followers.Users {
			followersMap[user.Username] = true
		}
	}

	var result []string
	following := user.Following()
	for following.Next() {
		for _, user := range following.Users {
			if _, exists := followersMap[user.Username]; !exists {
				result = append(result, user.Username)
			}
		}
	}

	return &result, nil
}

// GetUnfollowedUsersFile ...
func (i *Instagram) GetUnfollowedUsersFile(username, path string) error {
	users, err := i.GetUnfollowedUsers(username)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, user := range *users {
		f.WriteString(user + "\n")
	}

	return nil
}
