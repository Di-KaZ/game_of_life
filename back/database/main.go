package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb(file string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
	db.Create(&User{ID: 1, Password: "123456", Name: "23Alive", Width: 50, Height: 50, Alive: 20})
	db.Create(&User{ID: 2, Password: "123456", Name: "23Alive_big", Width: 100, Height: 100, Alive: 1000})
	return db
}
