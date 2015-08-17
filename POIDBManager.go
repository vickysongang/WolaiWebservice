package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type POIDBManager struct {
	dbClient *sql.DB
}


func NewPOIDBManager() POIDBManager {
	dbClient, _ := sql.Open("mysql",
		Config.Database.Username+":"+
			Config.Database.Password+"@"+
			Config.Database.Method+"("+
			Config.Database.Address+":"+
			Config.Database.Port+")/"+
			Config.Database.Database)
	return POIDBManager{dbClient: dbClient}
}

