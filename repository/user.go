package repository

import (
	"gorm.io/gorm"
	"sync"
)

type User struct {
	Id            int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name          string `gorm:"column:name;not null"`
	Password      string `gorm:"column:password;type:varchar(100);not null"`
	FollowCount   int64  `gorm:"column:follow_count;not null"`
	FollowerCount int64  `gorm:"column:follower_count;not null"`
	IsFollow      bool   `gorm:"column:is_follow;"`
}

type UserDAO struct {
}

var userDao *UserDAO // UserDAO这个类只有一个实例
var userOnce sync.Once

// 单例模式我的理解是：如果外界想要操作数据库，只能通过这个单例，这个单例有且仅会被创建一次

func NewUserDao() *UserDAO { // 单例设计模式：懒汉模式
	// 这个函数暴露唯一的接口对外使用
	userOnce.Do(
		func() {
			userDao = &UserDAO{}
		})
	return userDao
}

func (*UserDAO) AddUser(user *User) error {
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (*UserDAO) QueryByName(name string) (*User, error) {
	// 根据name查询，返回查询到的记录地址，此处函数的返回值是指针的目的就是返回User记录所在的地址
	var user User
	err := db.Where("name=?", name).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (*UserDAO) QueryById(id int64) (*User, error) {
	// 根据ID查询
	var user User
	err := db.Where("id=?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (*UserDAO) AddFollow(id int64) error {
	// 根据ID添加关注
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follow_count", gorm.Expr("follow_count+?", 1))
	if err != nil {
		panic(err)
	}
	return nil
}

func (*UserDAO) RmFollow(id int64) error {
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follow_count", gorm.Expr("follow_count-?", 1))
	if err != nil {
		panic(err)
	}
	return nil
}

func (*UserDAO) AddFollower(id int64) error {
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follower_count", gorm.Expr("follower_count+?", 1))
	if err != nil {
		panic(err)
	}
	return nil
}

func (*UserDAO) RmFollower(id int64) error {
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follower_count", gorm.Expr("follower_count-?", 1))
	if err != nil {
		panic(err)
	}
	return nil
}
