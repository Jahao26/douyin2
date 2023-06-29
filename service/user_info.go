package service

import (
	"douyin/repository"
)

func QueryUserById(uid int64) (*repository.User, error) {
	user, err := repository.NewUserDao().QueryById(uid)
	if err != nil {
		return &repository.User{}, err
	}
	return &repository.User{
		Id:            uid,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
		FavoriteCount: user.FavoriteCount,
		WorkCount:     user.WorkCount,
	}, nil
}
