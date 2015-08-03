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
