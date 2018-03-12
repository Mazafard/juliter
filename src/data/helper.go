package data

import (
	"db"
	"strings"
	"twitter"
)

func CheckTweet(username string, tweets []twitter.Tweet) {
	db.CreateTweets(username, tweets)

	quit := make(chan int)

	go extractTermsForUser(username, quit)

	for {
		select {
		case <-quit:
			return
		}
	}
}

func extractTermsForUser(username string, quit chan<- int) {
	termsDictionary := map[string]*twitter.TermDoc{}

	db.TweetsByUser(username, func(tweet *twitter.Tweet) {
		terms := strings.Split(normalizeText(tweet.Text), " ")

		for _, term := range terms {
			_, ok := termsDictionary[term]

			if ok {
				termsDictionary[term].Count++
			} else {
				termsDictionary[term] = &twitter.TermDoc{TweetId: tweet.Id, Term: term, Count: 1}
			}
		}
	})

	db.SaveTerms(username, termsDictionary)

	quit <- 1
}

func normalizeText(text string) string {
	return strings.ToLower(text)
}
