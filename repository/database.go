package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return err
	}
	err = db.AutoMigrate(&User{}, &Video{}, &Favorite{}, &Ralation{}, &Comment{})
	if err != nil {
		panic(err)
	}
	fmt.Println("****************************************")
	fmt.Println("MYSQL is connected")
	fmt.Println("****************************************")
	return err
}
