package pkg

import (
	"log"

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

// GetFollowers ...
func (i *Instagram) GetFollowers(username string) (*[]string, error) {
	user, err := i.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result []string
	followers := user.Followers()
	for followers.Next() {
		for _, user := range followers.Users {
			result = append(result, user.Username)
		}
	}

	return &result, nil
}

// GetFollowing ...
func (i *Instagram) GetFollowing(username string) (*[]string, error) {
	user, err := i.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result []string
	followers := user.Following()
	for followers.Next() {
		for _, user := range followers.Users {
			result = append(result, user.Username)
		}
	}

	return &result, nil
}

// GetUnfollowedUsers ...
func (i *Instagram) GetUnfollowedUsers(username string) (*[]string, error) {
	followers, err := i.GetFollowers(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	following, err := i.GetFollowing(username)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	followersMap := make(map[string]bool, 0)
	for _, user := range *followers {
		followersMap[user] = true
	}

	var result []string
	for _, user := range *following {
		if _, exists := followersMap[user]; !exists {
			result = append(result, user)
		}
	}

	return &result, nil
}
