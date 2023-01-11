package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Driver Connect DB
)

const (
	host     = "172.17.0.2"
	database = "devbook"
	user     = "root"
	password = "root"
	port     = "3306"
)

// Connect open connect DB
func Connect() (*sql.DB, error) {
	stringConnect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, database)
	db, erro := sql.Open("mysql", stringConnect)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
