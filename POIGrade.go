package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type POIGrade struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Pid  int64  `json:"pid"`
}

type POIGrades []POIGrade

func (dbm *POIDBManager) QueryGradeList() POIGrades {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, name, pid FROM grade`)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var id int64
	var name string
	var pid int64
	grades := make(POIGrades, 0)

	for rows.Next() {
		err = rows.Scan(&id, &name, &pid)

		grades = append(grades, POIGrade{Id: id, Name: name, Pid: pid})
	}

	return grades
}
