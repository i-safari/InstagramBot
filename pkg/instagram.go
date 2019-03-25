package pkg

import (
	"gopkg.in/ahmdrz/goinsta.v2"
)

// CreateInstagram ...
func CreateInstagram(username, password string) (*Instagram, error) {
	var instagram Instagram
	if err := instagram.SetAccount(username, password); err != nil {
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
		return err
	}

	return nil
}

// GetUserByUsername ...
func (i *Instagram) GetUserByUsername(username string) (*goinsta.User, error) {
	user, err := i.Account.Profiles.ByName(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUnfollowedUsers ...
func (i *Instagram) GetUnfollowedUsers(username string) (*[]goinsta.User, error) {
	user, err := i.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	followersMap := make(map[string]bool, 0)
	followers := user.Followers()
	for followers.Next() {
		for _, user := range followers.Users {
			followersMap[user.Username] = true
		}
	}

	var result []goinsta.User
	following := user.Following()
	for following.Next() {
		for _, user := range following.Users {
			if _, exists := followersMap[user.Username]; !exists {
				result = append(result, user)
			}
		}
	}

	return &result, nil
}
