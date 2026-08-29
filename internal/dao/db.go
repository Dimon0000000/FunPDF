package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global gorm database handle, initialized once at startup.
var DB *gorm.DB

// InitMysql connects to MySQL and stores the handle in the package-level DB.
func InitMysql(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

// InitSqlite connects to Sqlite and stores the handle in the package-level DB.
func InitSqlite(dsn string) error {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
