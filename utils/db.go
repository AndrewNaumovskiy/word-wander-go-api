package utils

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type DbUserModel struct {
	Id               int
	MongoId          string
	Email            string
	Password         string
	RegistrationDate string
}

type DicGetWordsResponseModel struct {
	Id          int                          `json:"id"`
	Word        string                       `json:"word"`
	Translation string                       `json:"translation"`
	Collections []DicGetWordsCollectionModel `json:"collections"`
}

type DicGetWordsCollectionModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DicGetWordsDbModel struct {
	Id             int
	Original       string
	Translation    string
	CollectionId   int
	CollectionName string
}

var db *sql.DB

func GetUserByEmail(email string) (*DbUserModel, error) {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user DbUserModel

	rows := db.QueryRow("SELECT id, mongoId, email, password, registrationDate FROM users_ww WHERE email = ?", email)

	if err := rows.Scan(&user.Id, &user.MongoId, &user.Email, &user.Password, &user.RegistrationDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error scanning user: %v", err)
		return nil, err
	}

	return &user, nil
}

func GetWordsByUserId(userId float64) ([]DicGetWordsResponseModel, error) {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var words []DicGetWordsDbModel

	rows, err := db.Query("SELECT w.Id AS WordId,w.Original,w.Translation,c.Id AS CollectionId,c.`Name` AS CollectionName FROM users_ww u JOIN collections c ON c.UserID = u.Id JOIN collections_words cw ON cw.CollectionId = c.Id JOIN words w ON w.Id = cw.WordId WHERE u.Id = ? ORDER BY w.UpdatedAt desc", userId)

	if err != nil {
		log.Printf("Error querying words: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var word DicGetWordsDbModel
		if err := rows.Scan(&word.Id, &word.Original, &word.Translation, &word.CollectionId, &word.CollectionName); err != nil {
			log.Printf("Error scanning word: %v", err)
			return nil, err
		}
		words = append(words, word)
	}

	var collectionsMap = make(map[int]DicGetWordsCollectionModel)

	for _, word := range words {
		if _, exists := collectionsMap[word.CollectionId]; !exists {
			collectionsMap[word.CollectionId] = DicGetWordsCollectionModel{
				Id:   word.CollectionId,
				Name: word.CollectionName,
			}
		}
	}

	var wordsResult []DicGetWordsResponseModel
	for _, word := range words {
		collection := collectionsMap[word.CollectionId]

		// check if the collection already exists in the response
		found := false
		for i, existingWord := range wordsResult {
			if existingWord.Id == word.Id {
				existingWord.Collections = append(existingWord.Collections, collection)
				wordsResult[i] = existingWord
				found = true
				break
			}
		}

		if found {
			continue
		}

		wordsResult = append(wordsResult, DicGetWordsResponseModel{
			Id:          word.Id,
			Word:        word.Original,
			Translation: word.Translation,
			Collections: []DicGetWordsCollectionModel{collection},
		})
	}

	return wordsResult, nil
}

type TraningWordsAmountModel struct {
	Id      int    `json:"id"`
	Traning string `json:"traning"`
	Amount  int    `json:"amountOfWords"`
}

func GetTrainingAmountWordsByUserId(userId float64) ([]TraningWordsAmountModel, error) {
	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var words []TraningWordsAmountModel

	row, err := db.Query(GetTrainingAmountWordsQuery(), userId)
	if err != nil {
		log.Printf("Error querying training words amount: %v", err)
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var word TraningWordsAmountModel
		if err := row.Scan(&word.Id, &word.Traning, &word.Amount); err != nil {
			log.Printf("Error scanning training words amount: %v", err)
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

func GetTrainingAmountWordsQuery() string {
	return `
		SELECT
			t.Id AS TrainingId,
			t.Name AS TrainingName,
			COUNT(tw.WordsId) AS WordCount
		FROM
			trainings t
		LEFT JOIN
			trainings_words tw ON tw.TrainingsID = t.Id
		WHERE
			t.UserId = ?
		GROUP BY
			t.Id, t.Name
		ORDER BY
			t.Id;
	`
}

func initDb() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = "the_guv"
	cfg.Passwd = "turn-table"
	cfg.Net = "tcp"
	cfg.Addr = "autoloadit-aurora-cluster.cluster-cfavjdqmfbhb.eu-west-1.rds.amazonaws.com"
	cfg.DBName = "autoloadit_testing"

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
