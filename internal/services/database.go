package services

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS posts(
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"slug" TEXT,
		"content" TEXT,
		"status" TEXT,
		"visible" BOOL,
		"created_at" DATETIME,
    	"updated_at" DATETIME

	);`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		return nil, err

	}

	statement.Exec()

	return db, nil
}
