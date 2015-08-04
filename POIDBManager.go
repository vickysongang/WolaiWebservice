package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type POIDBManager struct {
	dbClient *sql.DB
}

func NewPOIDBManager() POIDBManager {
	dbClient, _ := sql.Open("mysql", DB_URL)
	return POIDBManager{dbClient: dbClient}
}
