package data

import (
	"db"
	"fmt"
	"os"
	"twitter"
)

func ExampleCheckTweet() {
	os.Setenv("Juli-ENV", "test")

	db.Clear()

	tweets := []twitter.Tweet{
		twitter.Tweet{"hi I'm new!", 1},
		twitter.Tweet{"سلام من جدیدم", 2},
		twitter.Tweet{"خدافظ", 3},
		twitter.Tweet{"ایت یم تویت تست است", 4},
		twitter.Tweet{"قشنگه", 5},
	}

	CheckTweet("Mazafard", tweets)

	fmt.Printf("%d tweets inserted\n", db.TotalTweets())

	db.TermsByUser("chischaschos", func(termDoc *twitter.TermDoc) {
		fmt.Printf("%v\n", termDoc)
	})


}

func ExampleNormalizeText() {
	fmt.Println(normalizeText("تست"))

}
