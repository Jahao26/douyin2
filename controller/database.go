package controller

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Print("MYSQL is connected")
	err = db.AutoMigrate(&User{}, &UserLogin{}, &Relation{})
	if err != nil {
		panic(err)
	}
	//DB = db
	//return db
}
