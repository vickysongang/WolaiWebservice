package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const DB_URL_DEV = "poi:public11223@tcp(poianalytics.mysql.rds.aliyuncs.com:3306)/wolai"

type POIDBManager struct {
	dbClient *sql.DB
}

func NewPOIDBManager() POIDBManager {
	dbClient, _ := sql.Open("mysql", DB_URL_DEV)
	return POIDBManager{dbClient: dbClient}
}

func (dbm *POIDBManager) GetUserById(userId int) *POIUser {
	var nickname string
	var avatar string
	var gender int
	var accessRight int

	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT nickname, avatar, gender, access_right FROM users WHERE id = ?`)
	defer stmtQuery.Close()

	if err != nil {
		panic(err.Error())
	}

	rowUser := stmtQuery.QueryRow(userId)
	err = rowUser.Scan(&nickname, &avatar, &gender, &accessRight)

	user := NewPOIUser(userId, nickname, avatar, gender, accessRight)
	return &user
}

func (dbm *POIDBManager) GetUserByPhone(phone string) *POIUser {
	var userId int
	var nickname string
	var avatar string
	var gender int
	var accessRight int

	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, nickname, avatar, gender, access_right FROM users WHERE phone = ?`)
	defer stmtQuery.Close()

	if err != nil {
		panic(err.Error())
	}

	rowUser := stmtQuery.QueryRow(phone)
	err = rowUser.Scan(&userId, &nickname, &avatar, &gender, &accessRight)
	if err != nil {
		return nil
	}

	user := NewPOIUser(userId, nickname, avatar, gender, accessRight)
	return &user
}

func (dbm *POIDBManager) InsertUser(phone string) int {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO users(phone) VALUES(?)`)
	defer stmtInsert.Close()

	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	result, _ := stmtInsert.Exec(phone)
	id, _ := result.LastInsertId()
	return int(id)
}

func (dbm *POIDBManager) UpdateUserInfo(userId int, nickname string, avatar string, gender int) *POIUser {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE users SET nickname = ?, avatar = ?, gender = ? WHERE id = ?`)
	defer stmtUpdate.Close()

	if err != nil {
		panic(err.Error())
	}

	_, err = stmtUpdate.Exec(nickname, avatar, gender, userId)

	user := NewPOIUser(userId, nickname, avatar, gender, 3)
	return &user
}
