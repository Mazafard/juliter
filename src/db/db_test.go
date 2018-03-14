package db

import (
	"fmt"
	"os"
	"twitter"
)

func ExampleTweetsByUser() {
	os.Setenv("JULITER-ENV", "test")

	Clear()

	tweets := []twitter.Tweet{
		twitter.Tweet{"سلام سلام", 226233195904875998},
		twitter.Tweet{"خوب حالم", 3765231951076753242},
	}

	CreateTweets("mazafard", tweets)

	TweetsByUser("chischaschos", func(tweet *twitter.Tweet) {
		fmt.Printf("%v\n", tweet)
	})


}
