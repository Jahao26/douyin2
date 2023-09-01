package repository

import (
	"gorm.io/gorm"
	"sync"
)

type User struct {
	Id            int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" redis:"id"`
	Name          string `gorm:"column:name;not null" redis:"name"`
	Password      string `gorm:"column:password;type:varchar(100);not null" redis:"password"`
	FollowCount   int64  `gorm:"column:follow_count;not null" redis:"follow_count"`
	FollowerCount int64  `gorm:"column:follower_count;not null" redis:"follower_count"`
	IsFollow      bool   `gorm:"column:is_follow;" redis:"is_follow"`
	FavoriteCount int64  `gorm:"column:favorite_count;not null" redis:"favorite_count"`
	WorkCount     int64  `gorm:"column:work_count;not null" redis:"work_count"`
}

type UserDAO struct {
}

var userDao *UserDAO // UserDAO这个类只有一个实例
var userOnce sync.Once

// 单例模式：如果外界想要操作数据库，只能通过这个单例，这个单例有且仅会被创建一次

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

// 根据ID增加关注数量
func (*UserDAO) AddFollow(id int64) error {

	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follow_count", gorm.Expr("follow_count+?", 1)).Error
	if err != nil {
		//panic(err)
		return err
	}
	return nil
}

func (*UserDAO) RmFollow(id int64) error {
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follow_count", gorm.Expr("follow_count-?", 1)).Error
	if err != nil {
		//panic(err)
		return err
	}
	return nil
}

// 根据ID增加粉丝数量
func (*UserDAO) AddFollower(id int64) error {
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follower_count", gorm.Expr("follower_count+?", 1)).Error
	if err != nil {
		//panic(err)
		return err
	}
	return nil
}

func (*UserDAO) RmFollower(id int64) error {
	err := db.Model(&User{}).Where("id=?", id).UpdateColumn("follower_count", gorm.Expr("follower_count-?", 1)).Error
	if err != nil {
		//panic(err)
		return err
	}
	return nil
}

func (*UserDAO) AddworkCount(uid int64) error {
	db.Model(&User{}).Where("id=?", uid).UpdateColumn("work_count", gorm.Expr("work_count+?", 1))
	return nil
}

func (*UserDAO) AddfavoriteCount(uid int64) error {
	db.Model(&User{}).Where("id=?", uid).UpdateColumn("favorite_count", gorm.Expr("favorite_count+?", 1))
	return nil
}

func (*UserDAO) RmworkCount(uid int64) error {
	db.Model(&User{}).Where("id=?", uid).UpdateColumn("work_count", gorm.Expr("work_count-?", 1))
	return nil
}

func (*UserDAO) RmfavoriteCount(uid int64) error {
	db.Model(&User{}).Where("id=?", uid).UpdateColumn("favorite_count", gorm.Expr("favorite_count-?", 1))
	return nil
}
