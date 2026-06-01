package handlers

import (
	"database/sql"

	"gorm.io/gorm"
)

var db *sql.DB
var gormDB *gorm.DB

func SetDB(database *sql.DB) {
	db = database
}

func SetGormDB(database *gorm.DB) {
	gormDB = database
}
