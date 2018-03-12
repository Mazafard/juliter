package main

import (
	"data"
	"db"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	tw "twitter"
)

var usernames string

func init() {
	flag.StringVar(&usernames, "u", "Mazafard,shib,baam", "the users who play with them")

	db.Settings().StorePath = settingsPath()
}

func settingsPath() string {
	currentUser, userError := user.Current()

	if userError != nil {
		panic(userError)
	}

	settingsFile := currentUser.HomeDir + "/.juli.json"
	settingsPath, filepathError := filepath.Abs(settingsFile)

	if filepathError != nil {
		panic(filepathError)
	}

	return settingsPath
}

func authValues() (string, string) {
	consumerKey := db.Settings().Get("consumer-key")

	if consumerKey == "" {
		consumerKey := os.Getenv("CONSUMER_KEY")

		if consumerKey == "" {
			panic("not exist CONSUMER_KEY")
		} else {
			db.Settings().Set("consumer-key", consumerKey)
		}
	}

	consumerSecret := db.Settings().Get("consumer-secret")

	if consumerSecret == "" {
		consumerSecret := os.Getenv("CONSUMER_SECRET")

		if consumerSecret == "" {
			panic("not exist CONSUMER_SECRET")
		} else {
			db.Settings().Set("consumer-secret", consumerSecret)
		}
	}

	return consumerKey, consumerSecret
}

func main() {
	flag.Parse()

	if usernames == "Mazafard,shib,baam" {
		flag.PrintDefaults()
	} else {
		consumerKey, consumerSecret := authValues()
		fmt.Println(consumerKey)
		fmt.Println(consumerSecret)
		twitter := tw.New(consumerKey, consumerSecret)

		tweetsChannel := make(chan *tw.User)
		splitUsernames := strings.Split(usernames, ",")
		usersCount := len(splitUsernames)
		messagesCount := 0

		for _, username := range splitUsernames {
			go twitter.FetchTweetsOf(username, tweetsChannel)
		}

	main:
		for {
			select {
			case twitterUser := <-tweetsChannel:
				messagesCount++
				data.CheckTweet(twitterUser.Name, twitterUser.Tweets)

				if messagesCount == usersCount {
					break main
				}
			}
		}

		for _, username := range splitUsernames {
			db.TermsByUser(username, func(termDoc *tw.TermDoc) {
				fmt.Printf("%v\n", termDoc)
			})
		}
	}
}
