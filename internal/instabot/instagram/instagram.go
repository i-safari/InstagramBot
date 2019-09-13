package instagram

import "github.com/ahmdrz/goinsta/v2"

// CreateInstagram ...
func CreateInstagram(username, password string) (*Instagram, error) {
	var instagram Instagram
	if err := instagram.setAccount(username, password); err != nil {
		return nil, err
	}

	return &instagram, nil
}

// Instagram ...
type Instagram struct {
	Account *goinsta.Instagram
}

// SetAccount ...
func (i *Instagram) setAccount(username, password string) error {
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
