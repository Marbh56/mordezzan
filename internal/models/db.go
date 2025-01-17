package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Models struct {
    // We'll add more models here as we create them
    Users UserModel
}

func NewModels(db *sql.DB) Models {
    return Models{
        Users: UserModel{DB: db},
    }
}
