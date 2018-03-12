package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mxk/go-sqlite/sqlite3"
	"os"
	"os/user"
	"path/filepath"
	"twitter"
)

func DB() *sql.DB {
	db, connErr := sql.Open("sqlite3", dbName())

	if connErr != nil {
		panic(connErr)
	}

	setupStatements := []string{
		"CREATE TABLE IF NOT EXISTS tweets(id, text, username, source)",
		"CREATE TABLE IF NOT EXISTS tweet_terms(username, tweet_id, term, count)",
		"CREATE UNIQUE INDEX IF NOT EXISTS unique_tweet_idx ON tweets(id, username)",
	}

	for _, stmt := range setupStatements {
		_, stmtError := db.Exec(stmt)

		if stmtError != nil {
			panic(stmtError)
		}
	}

	return db
}

func TweetsByUser(username string, callback func(*twitter.Tweet)) {
	db := DB()
	defer db.Close()

	rows, selectError := db.Query("SELECT id, text, source FROM tweets WHERE username = ?", username)

	if selectError != nil {
		panic(selectError)
	}

	for rows.Next() {
		var tweet twitter.Tweet

		rowError := rows.Scan(&tweet.Id, &tweet.Text, &tweet.Source)

		if rowError != nil {
			fmt.Printf("%#v\n", rowError)
		} else {
			callback(&tweet)
		}
	}
}

func TermsByUser(username string, callback func(*twitter.TermDoc)) {
	db := DB()
	defer db.Close()

	rows, selectError := db.Query("SELECT tweet_id, term, count FROM tweet_terms WHERE username = ?", username)

	if selectError != nil {
		panic(selectError)
	}

	for rows.Next() {
		var termDoc twitter.TermDoc
		rowError := rows.Scan(&termDoc.TweetId, &termDoc.Term, &termDoc.Count)

		if rowError != nil {
			fmt.Println(rowError)
		} else {
			callback(&termDoc)
		}
	}
}

func CreateTweets(username string, tweets []twitter.Tweet) {
	db := DB()
	defer db.Close()

	stmt, prepareError := db.Prepare("INSERT OR REPLACE INTO tweets VALUES(?, ?, ?, ?)")

	if prepareError != nil {
		panic(prepareError)
	}

	for _, tweet := range tweets {
		_, execError := stmt.Exec(tweet.Id, tweet.Text, username, tweet.Source)

		if execError != nil {
			fmt.Println(execError)
		}

	}
}

func SaveTerms(username string, termsDictionary map[string]*twitter.TermDoc) {
	db := DB()
	defer db.Close()

	insertStmt, stmtError := db.Prepare("INSERT INTO tweet_terms (username, tweet_id, term, count) VALUES (?, ?, ?, ?)")

	if stmtError != nil {
		panic(stmtError)
	}

	for _, termDoc := range termsDictionary {
		_, ie := insertStmt.Exec(username, termDoc.TweetId, termDoc.Term, termDoc.Count)

		if ie != nil {
			fmt.Println(ie)
		}
	}
}

func Clear() {
	db := DB()
	defer db.Close()

	cleanUpStmts := []string{
		"DELETE FROM tweets",
		"DELETE FROM tweet_terms",
		"DROP INDEX unique_tweet_idx",
	}

	for _, stmt := range cleanUpStmts {
		_, stmtError := db.Exec(stmt)

		if stmtError != nil {
			panic(stmtError)
		}
	}
}

func TotalTweets() int {
	db := DB()
	defer db.Close()

	stmt, countError := db.Query("SELECT count(*) from tweets")

	if countError != nil {
		panic(countError)
	}

	var count int

	stmt.Next()
	stmt.Scan(&count)

	return count
}

func dbName() string {
	var dbName string

	if os.Getenv("Juli-ENV") == "test" {
		dbName = "../juli-db-test"

	} else {
		currentUser, userError := user.Current()

		if userError != nil {
			panic(userError)
		}

		dbName = currentUser.HomeDir + "/.juli-db"
	}

	path, filepathError := filepath.Abs(dbName)

	if filepathError != nil {
		panic(filepathError)
	}

	return path
}
