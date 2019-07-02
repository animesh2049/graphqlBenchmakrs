package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	dbname   = "test"
)

func topKQuestions(db *sql.DB) {
	sqlQry := fmt.Sprintf("select p1.Id,p1.Score as qscore, count(p2.Id) as " +
		"answercount,sum(p2.Score), p1.score*p1.answercount+sum(p2.score) as finalscore from " +
		"posts as p1 inner join posts as p2 on p1.Id=p2.ParentId group by p1.Id  order by " +
		"finalscore desc;")

	var (
		id          int
		qscore      int
		answerCount int
		sum         int
		finalscore  int
	)

	result, err := db.Query(sqlQry)
	if err != nil {
		panic(err)
	}

	for result.Next() {
		err := result.Scan(&id, &qscore, &answerCount, &sum, &finalscore)
		if err != nil {
			panic(err)
		}

		// fmt.Println(id, qscore, answerCount, sum, finalscore)
	}
}

func questionAnswer(db *sql.DB, quesID int) {
	var (
		id                    int
		score                 int
		text                  string
		displayName           string
		lastEditorDisplayName string
		creationDate          string
		lastActivityDate      string
	)

	answersSQL := fmt.Sprintf("select p.Id, p.score, p.lastEditorDisplayName, p.creationDate, "+
		"p.lastActivityDate, u.DisplayName from posts as p, users as u where p.parentId=%d and "+
		"p.ownerUserId=u.Id order by p.score desc;", quesID)

	results, err := db.Query(answersSQL)
	if err != nil {
		panic(err)
	}

	ids := make([]int, 0)
	for results.Next() {
		err = results.Scan(&id, &score, &lastEditorDisplayName, &creationDate,
			&lastActivityDate, &displayName)
		ids = append(ids, id)
	}

	for _, val := range ids {
		answerTextSQL := fmt.Sprintf("select text from posthistory where postId=%d order by "+
			"creationDate desc limit 1;", val)

		results, err := db.Query(answerTextSQL)
		if err != nil {
			panic(err)
		}

		for results.Next() {
			err = results.Scan(&text)
			if err != nil {
				panic(err)
			}
		}

		answerComments := fmt.Sprintf("select c.text, c.score, c.creationDate, u.displayName from "+
			"comments as c, users as u where c.postId=%d and c.userId=u.Id order by c.score desc;",
			val)

		results, sqlErr := db.Query(answerComments)
		if sqlErr != nil {
			panic(sqlErr)
		}

		for results.Next() {
			sqlErr = results.Scan(&text, &score, &creationDate, &displayName)
			if sqlErr != nil {
				panic(sqlErr)
			}
		}
	}
}

func questions(db *sql.DB, quesID int) {
	var (
		score        int
		answerCount  int
		creationDate string
		displayName  string
	)

	questionSQL := fmt.Sprintf("select p1.Score, count(p2.Id) as answerCount, p1.CreationDate, "+
		"u.DisplayName from posts as p1 inner join posts as p2 on p1.Id=p2.ParentId, users as u "+
		"where p1.Id=%d and p1.OwnerUserId=u.Id group by p1.Id, u.DisplayName", quesID)

	results, err := db.Query(questionSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&score, &answerCount, &creationDate, &displayName)
		if err != nil {
			panic(err)
		}
	}

	var text string
	qstnBodySQL := fmt.Sprintf("select text from posthistory where postId=%d order by "+
		"creationDate desc limit 1;", quesID)

	results, err = db.Query(qstnBodySQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&text)
		if err != nil {
			panic(err)
		}
	}

	var commentCount int
	commentCountSQL := fmt.Sprintf("select count(Id) from comments where postId=%d", quesID)

	results, err = db.Query(commentCountSQL)
	if err != nil {
		panic(err)
	}

	for results.Next() {
		err = results.Scan(&commentCount)
		if err != nil {
			panic(err)
		}
	}

	commentText := fmt.Sprintf("select c.text, c.score, c.creationDate, u.displayName from "+
		"comments as c, users as u where c.postId=%d and c.userId=u.Id;", quesID)

	results, sqlErr := db.Query(commentText)
	if sqlErr != nil {
		panic(sqlErr)
	}

	for results.Next() {
		sqlErr = results.Scan(&text, &score, &creationDate, &displayName)
		if sqlErr != nil {
			panic(sqlErr)
		}
	}
}

func executeSQL(db *sql.DB) {

	postID := 1
	sqlStatements := make([]string, 4)

	sqlStatements[0] = fmt.Sprintf("select p1.Id,p1.Score as qscore, count(p2.Id) as " +
		"answercount,sum(p2.Score), p1.score*p1.answercount+sum(p2.score) as finalscore from " +
		"posts as p1 inner join posts as p2 on p1.Id=p2.ParentId group by p1.Id  order by " +
		"finalscore desc;")

	sqlStatements[1] = fmt.Sprintf("select p.Score, p.ViewCount, p.Body, p.Title, p.Tags, p.AnswerCount, "+
		"p.CommentCount, u.CreationDate, u.DisplayName from posts as p, users as u where "+
		"p.Id=%d and u.Id=p.OwnerUserId;", postID)

	sqlStatements[2] = fmt.Sprintf("select p.Score, p.ViewCount, p.Body, p.Title, p.Tags, "+
		"p.LastEditorDisplayName, p.LastActivityDate from posts as p where p.ParentId=%d "+
		"order by p.Score desc;", postID)

	sqlStatements[3] = fmt.Sprintf("```select p.Id, c.TextId, c.Score, c.CreationDate, "+
		"u.DisplayName from posts as p, comments as c, users as u where p.ParentId=%d "+
		"and c.PostId=p.Id and c.UserId=u.Id;", postID)

	topKQuestions(db)
	questions(db, 1)
	questionAnswer(db, 1)

}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	executeSQL(db)
}

