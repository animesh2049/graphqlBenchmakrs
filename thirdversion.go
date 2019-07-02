package main

import (
	"database/sql"
	"fmt"
	"strings"
)

func getAnswersCommentsUsers(db *sql.DB, postID int) {
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

	var answerIDs, userIDs []string
	for results.Next() {
		err = results.Scan(&id, &creationDate, &score, &viewCount, &body, &ownerUserID, &title, &tags)
		if err != nil {
			panic(err)
		}

		answerIDs = append(answerIDs, string(id))
		userIDs = append(userIDs, string(ownerUserID))

		// getComments(db, id)
		// getOwner(db, id)
	}

	var (
		text        string
		userID      int
		reputation  string
		displayName string
	)

	commentSQL := "select c.text, c.score, c.creationDate, c.userId from comments as c where c.postId in (" + strings.Join(answerIDs, ",") + ");"

	results, err = db.Query(commentSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		if err = results.Scan(&text, &score, &creationDate, &userID); err != nil {
			panic(err)
		}
	}

	userSQL := "select u.reputation, u.displayName, from users as u where u.id in (" + strings.Join(userIDs, ",") + ");"

	results, err = db.Query(userSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		if err = results.Scan(&reputation, &displayName); err != nil {
			panic(err)
		}
	}
}

func getPostData2(db *sql.DB, postID int) {
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

		getAnswersCommentsUsers(db, postID)
	}
}

