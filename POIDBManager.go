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

func (dbm *POIDBManager) GetUserById(userId int64) *POIUser {
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var accessRight int64

	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT nickname, avatar, gender, access_right FROM users WHERE id = ?`)
	defer stmtQuery.Close()

	if err != nil {
		panic(err.Error())
	}

	rowUser := stmtQuery.QueryRow(userId)
	err = rowUser.Scan(&nicknameNS, &avatarNS, &gender, &accessRight)

	if err == sql.ErrNoRows {
		return nil
	}

	nickname := ""
	if nicknameNS.Valid {
		nickname = nicknameNS.String
	}

	avatar := ""
	if avatarNS.Valid {
		avatar = avatarNS.String
	}

	user := NewPOIUser(int64(userId), nickname, avatar, gender, accessRight)
	return &user
}

func (dbm *POIDBManager) GetUserByPhone(phone string) *POIUser {
	var userId int64
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var accessRight int64

	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, nickname, avatar, gender, access_right FROM users WHERE phone = ?`)
	defer stmtQuery.Close()

	if err != nil {
		panic(err.Error())
	}

	rowUser := stmtQuery.QueryRow(phone)
	err = rowUser.Scan(&userId, &nicknameNS, &avatarNS, &gender, &accessRight)
	if err == sql.ErrNoRows {
		return nil
	}

	nickname := ""
	if nicknameNS.Valid {
		nickname = nicknameNS.String
	}

	avatar := ""
	if avatarNS.Valid {
		avatar = avatarNS.String
	}

	user := NewPOIUser(userId, nickname, avatar, gender, accessRight)
	return &user
}

func (dbm *POIDBManager) InsertUser(phone string) int64 {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO users(phone) VALUES(?)`)
	defer stmtInsert.Close()
	if err != nil {
		panic(err.Error())
	}

	result, err := stmtInsert.Exec(phone)
	if err != nil {
		panic(err.Error())
	}

	id, _ := result.LastInsertId()

	return id
}

func (dbm *POIDBManager) UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) *POIUser {
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
