package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type POIUser struct {
	UserId      int64  `json:"userId"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"gender"`
	AccessRight int64  `json:"accessRight"`
}

type POIUsers []POIUser

func NewPOIUser(userId int64, nickname string, avatar string, gender int64, accessRight int64) POIUser {
	user := POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender, AccessRight: accessRight}
	return user
}

func (dbm *POIDBManager) InsertUser(phone string) int64 {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO users(phone) VALUES(?)`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	result, err := stmtInsert.Exec(phone)
	if err != nil {
		panic(err.Error())
	}

	id, _ := result.LastInsertId()

	return id
}

func (dbm *POIDBManager) QueryUserById(userId int64) *POIUser {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT nickname, avatar, gender, access_right FROM users WHERE id = ?`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var accessRight int64

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

func (dbm *POIDBManager) QueryUserByPhone(phone string) *POIUser {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, nickname, avatar, gender, access_right FROM users WHERE phone = ?`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	var userId int64
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var accessRight int64

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

func (dbm *POIDBManager) UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) *POIUser {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE users SET nickname = ?, avatar = ?, gender = ? WHERE id = ?`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtUpdate.Close()

	_, err = stmtUpdate.Exec(nickname, avatar, gender, userId)
	if err != nil {
		panic(err.Error())
	}

	user := NewPOIUser(userId, nickname, avatar, gender, 3)

	return &user
}

func (dbm *POIDBManager) InsertUserOauth(userId int64, qqOpenId string) {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO user_oauth(user_id, open_id_qq) VALUES(?, ?)`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	_, err = stmtInsert.Exec(userId, qqOpenId)
	if err != nil {
		panic(err.Error())
	}
}

func (dbm *POIDBManager) QueryUserByQQOpenId(qqOpenId string) int64 {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT user_id FROM user_oauth WHERE open_id_qq = ?`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	var userId int64

	userRow := stmtQuery.QueryRow(qqOpenId)
	err = userRow.Scan(&userId)
	if err == sql.ErrNoRows {
		return -1
	}

	return userId
}
