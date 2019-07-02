package main

import (
	"database/sql"
	"fmt"
)

func getOwner(db *sql.DB, userID int) {
	var (
		reputation   string
		creationDate string
		displayName  string
		aboutMe      string
		views        int
		upvotes      int
		downvotes    int
	)

	userSQL := fmt.Sprintf("select u.reputation, u.creationDate, u.displayName, u.aboutMe, "+
		"u.views, u.upvotes, u.downvotes from users as u where u.id=%d", userID)

	results, err := db.Query(userSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&reputation, &creationDate, &displayName, &aboutMe, &views,
			&upvotes, &downvotes)
		if err != nil {
			panic(err)
		}
	}
}

func getComments(db *sql.DB, postID int) {
	var (
		text         string
		userID       int
		score        int
		creationDate string
	)

	commentSQL := fmt.Sprintf("select c.text, c.userID, c.score, c.creationDate from comments as "+
		"c where c.postID=%d", postID)

	results, err := db.Query(commentSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&text, &userID, &score, &creationDate)
		if err != nil {
			panic(err)
		}
	}
}

func getAnswers(db *sql.DB, postID int) {
	var (
		id           int
		creationDate string
		score        int
		viewCount    int
		body         string
		ownerUserID  int
		title        string
		tags         string
	)
	// Should I get answercount with a separate query or just count number of answers ?

	answerSQL := fmt.Sprintf("select p.id, p.creationDate, p.score, p.viewCount, p.body, p.ownerUserID, "+
		"p.title, p.tags from posts as p where p.parentId=%d order by p.score desc", postID)

	results, err := db.Query(answerSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&id, &creationDate, &score, &viewCount, &body, &ownerUserID, &title, &tags)
		if err != nil {
			panic(err)
		}

		getComments(db, id)
		getOwner(db, id)
	}
}

func getPostData(db *sql.DB, postID int) {
	var (
		id           int
		creationDate string
		score        int
		viewCount    int
		body         string
		tags         string
		ownerUserID  int
	)

	postSQL := fmt.Sprintf("select p.id, p.creationDate, p.score, p.viewCount, p.body, p.tags, "+
		"p.ownerUserId from posts as p where p.id=%d", postID)

	results, err := db.Query(postSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&id, &creationDate, &score, &viewCount, &body, &tags, &ownerUserID)
		if err != nil {
			panic(err)
		}

		getOwner(db, ownerUserID)
		getComments(db, postID)
		getAnswers(db, postID)
	}
}

