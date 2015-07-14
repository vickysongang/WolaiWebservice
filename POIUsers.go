package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type POIUser struct {
	UserId   int    `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
}

type POIUsers []POIUser

func NewPOIUser(userId int, nickname string, avatar string, gender int) POIUser {
	user := POIUser{userId, nickname, avatar, gender}
	return user
}

const DB_URL_DEV = "poi:public11223@tcp(poianalytics.mysql.rds.aliyuncs.com:3306)/wolai"

func POIUserLogin(phone string) (int, string) {
	dbClient, err := sql.Open("mysql", DB_URL_DEV)
	if err != nil {
		//TODO: Replace by proper error handling
		//panic(err.Error())
		fmt.Println(err.Error())
	}
	defer dbClient.Close()

	stmtQuery, err := dbClient.Prepare(
		`SELECT id, nickname, avatar, gender FROM users WHERE phone = ?`)
	if err != nil {
		//TODO: Replace by proper error handling
		//panic(err.Error())
		fmt.Println(err.Error())
	}
	defer stmtQuery.Close()

	rowsUser, err := stmtQuery.Query(phone)
	if err != nil {
		//panic(err.Error())
		fmt.Println(err.Error())
	}

	for rowsUser.Next() {
		var id int
		var nickname string
		var avatar string
		var gender int

		err = rowsUser.Scan(&id, &nickname, &avatar, &gender)
		if err == sql.ErrNoRows {
			return 2, "no result"
		}

		content, _ := json.Marshal(NewPOIUser(id, nickname, avatar, gender))
		return 0, string(content)
	}

	stmtInsert, err := dbClient.Prepare(
		`INSERT INTO users(phone) VALUES(?)`)
	if err != nil {
		//TODO: Replace by proper error handling
		//panic(err.Error())
		fmt.Println(err.Error())
	}
	defer stmtInsert.Close()

	result, _ := stmtInsert.Exec(phone)
	id2, _ := result.LastInsertId()
	content2, _ := json.Marshal(NewPOIUser(int(id2), "", "", 0))
	return 1001, string(content2)
}

func POIUserUpdateProfile(userId int, nickname string, avatar string, gender int) (int, string) {
	dbClient, err := sql.Open("mysql", DB_URL_DEV)
	if err != nil {
		//TODO: Replace by proper error handling
		//panic(err.Error())
		fmt.Println(err.Error())
	}
	defer dbClient.Close()

	stmtUpdate, err := dbClient.Prepare(
		`UPDATE users SET nickname = ?, avatar = ?, gender = ? WHERE id = ?`)
	if err != nil {
		//TODO: Replace by proper error handling
		//panic(err.Error())
		fmt.Println(err.Error())
	}

	_, err = stmtUpdate.Exec(nickname, avatar, gender, userId)

	content, _ := json.Marshal(NewPOIUser(userId, nickname, avatar, gender))
	return 0, string(content)
}
