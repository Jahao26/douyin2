package service

import (
	"douyin/repository"
)

//type Userinfo struct {
//	Id            int64
//	Name          string
//	FollowCount   int64
//	FollowerCount int64
//	IsFollow      bool
//}

func QueryUserById(uid int64) (*repository.User, error) {
	user, err := repository.NewUserDao().QueryById(uid)
	if err != nil {
		return &repository.User{}, err
	}
	//fmt.Println("*******************************************")
	//fmt.Println("debug name: ", user.Name)
	//fmt.Println("*******************************************")
	return &repository.User{
		Id:            uid,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}, nil
}
