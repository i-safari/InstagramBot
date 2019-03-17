package main

import (
	"fmt"
	"log"
	"os"

	"github.com/InstaSherlock/pkg"
)

func main() {
	instagram, err := pkg.CreateInstagram("_i_am_not_bot_", "Neversleeps5")
	if err != nil {
		log.Fatal(err)
	}

	unfollowedUsers, err := instagram.GetUnfollowedUsers("iness_m_")
	if err != nil {
		log.Fatal(err)
	}

	f, _ := os.Create("unfollowed_users.txt")

	for i, user := range *unfollowedUsers {
		_, _ = f.WriteString(fmt.Sprintf("%d - %s\n", i, user))
	}
}
