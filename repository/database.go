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
	fmt.Print("MYSQL is connected")
	return err
}
