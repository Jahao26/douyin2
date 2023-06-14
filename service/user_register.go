package service

import (
	"douyin/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type registerFlow struct {
	name     string
	password string
}

func Register(name string, password string) (int64, error) {
	// 外部接口 Register
	return NewRegisterFlow(name, password).Do() // Do?
}

func NewRegisterFlow(name string, password string) *registerFlow {
	return &registerFlow{name: name, password: password} // 创建新的类接口
}

func (f *registerFlow) Do() (int64, error) {
	id, err := f.userRegister()
	if err != nil {
		panic(err)
	}
	return id, nil
}

func (f *registerFlow) userRegister() (int64, error) {
	// 先检查用户是否存在
	if user, err := repository.NewUserDao().QueryByName(f.name); err != nil {
		return 0, err
	} else if err == nil && user.Id != 0 {
		return user.Id, errors.New("用户已存在")
	}

	// 新建用户
	hasspassword, err := bcrypt.GenerateFromPassword([]byte(f.password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	newUser := &repository.User{
		Id:            0,
		Name:          f.name,
		Password:      string(hasspassword),
		FollowCount:   0,
		FollowerCount: 0,
	}
	if err := repository.NewUserDao().AddUser(newUser); err != nil {
		return 0, err
	}
	return newUser.Id, nil
}
