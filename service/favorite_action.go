package service

import (
	"douyin/repository"
)

type favoriteFlow struct {
	uid         int64
	vid         int64
	action_type string
}

func FavoriteAction(uid int64, vid int64, action string) error {
	return NewFavoriteFlow(uid, vid, action).Do()
}

func NewFavoriteFlow(uid int64, vid int64, action string) *favoriteFlow {
	return &favoriteFlow{uid: uid, vid: vid, action_type: action}
}

func (f *favoriteFlow) Do() error {
	err := f.favoriteAction()
	if err != nil {
		return err
	}
	return nil
}

func (f *favoriteFlow) favoriteAction() error {
	if f.action_type == "1" { // 如果操作类型是 点赞
		// video中的赞数量+1
		err := repository.NewVideoDao().AddVideoFavorite(f.uid, f.vid) // 通过视频查询到video类
		if err != nil {
			return err
		}
		err = repository.NewFavoriteDao().AddFavorite(f.uid, f.vid)
		if err != nil {
			return err
		}
		// 用户喜欢数量加1
		err = repository.NewUserDao().AddfavoriteCount(f.uid)
		if err != nil {
			return err
		}

	} else { // 如果操作类型是 取消点赞
		err := repository.NewVideoDao().RmVideoFavorite(f.uid, f.vid) // 通过视频查询到video类
		if err != nil {
			return err
		}
		err = repository.NewFavoriteDao().RmFavorite(f.uid, f.vid)
		if err != nil {
			return err
		}
		err = repository.NewUserDao().RmfavoriteCount(f.uid)
		if err != nil {
			return err
		}
	}
	return nil
}
