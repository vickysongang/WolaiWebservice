package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type POISubject struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type POISubjects []POISubject

func (dbm *POIDBManager) QuerySubjectList() POISubjects {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, name FROM subject`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	rows, err := stmtQuery.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var id int64
	var name string
	subjects := make(POISubjects, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name)

		subjects = append(subjects, POISubject{Id: id, Name: name})
	}

	return subjects
}
