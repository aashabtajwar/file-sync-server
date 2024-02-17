package database

import (
	"database/sql"

	"github.com/aashabtajwar/th-server/errorhandling"
)

var databases []*sql.DB

func DBInit() *sql.DB {
	if len(databases) > 0 {
		return databases[0]
	}
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	errorhandling.DbConnectionError(err)
	databases = append(databases, db)
	return db
}
