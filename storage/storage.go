package storage

import (
	"database/sql"
)

var db *sql.DB

func Init() error {
	var err error
	db, err = sql.Open("mysql", name+":"+psw+"@/"+dbname)
	return err
}

func NewQuery(query string) (*sql.Rows, error) {
	rows, err := db.Query(query)
	return rows, err
}

func NewExec(exec string) (sql.Result, error) {
	result, err := db.Exec(exec)
	return result, err
}
