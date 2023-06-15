package service

import (
	"douyin/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type loginFlow struct {
	name     string
	password string
}

func UserLogin(name string, password string) (int64, error) {
	return NewLoginFlow(name, password).Do()
}

func NewLoginFlow(name string, password string) *loginFlow {
	return &loginFlow{name: name, password: password}
}

func (f *loginFlow) Do() (int64, error) {
	uid, err := f.userLogin()
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return uid, nil
}

func (f *loginFlow) userLogin() (int64, error) {
	// 检查用户是否存在
	user, err := repository.NewUserDao().QueryByName(f.name)
	if err != nil || user.Id == 0 {
		return 0, errors.New("User not exist")
	}
	// 验证用户密码

	hashpassword := []byte(f.password)  //登录输入的明文密码
	mypassword := []byte(user.Password) //数据库存储的加密哈希

	err = bcrypt.CompareHashAndPassword(mypassword, hashpassword)
	if err != nil {
		return 0, errors.New("Password error: " + err.Error())
	}

	// 返回验证成功的用户uid
	return user.Id, nil
}
