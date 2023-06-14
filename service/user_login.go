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
		return 0, errors.New("Error in login")
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
	hashpassword := []byte(f.password)
	mypassword := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashpassword, mypassword)
	if err != nil {
		return 0, errors.New("Password error")
	}

	// 返回验证成功的用户uid
	return user.Id, nil
}