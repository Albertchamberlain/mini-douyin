package conf

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database(connString string) {
	if db, err := gorm.Open(mysql.Open(connString), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(20)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Second * 30)
		DB = db
		migration()
	}
}
