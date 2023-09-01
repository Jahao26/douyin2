package service

import (
	"douyin/repository"
	"errors"
)

type UserInfoPage struct {
	// 个人主页情况
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	WorkCount     int64  `json:"work_count,omitempty"`
}

type InfoFlow struct {
	uid int64
}

func UserInfo(uid int64) (*UserInfoPage, error) {
	return NewInfoFlow(uid).Do()
}

func NewInfoFlow(uid int64) *InfoFlow {
	return &InfoFlow{uid: uid}
}

func (f *InfoFlow) Do() (*UserInfoPage, error) {
	user, err := f.userInfo()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return user, nil
}

func (f *InfoFlow) userInfo() (*UserInfoPage, error) {
	// 此处需要用redis更新数据库的followcount和followercount
	user, err := repository.NewUserDao().QueryById(f.uid)
	if err != nil {
		return nil, err
	}
	newInfoPage := UserInfoPage{
		Id:            user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
		FavoriteCount: user.FavoriteCount,
		WorkCount:     user.WorkCount,
	}
	return &newInfoPage, nil
}
