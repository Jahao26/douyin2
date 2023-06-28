package service

import "douyin/repository"

type favoriteFlow struct {
	uid         int64
	vid         int64
	action_type string
}

func FavoriteAction(uid int64, vid int64, action string) (int64, error) {
	return NewFavoriteFlow(uid, vid, action).Do()
}

func NewFavoriteFlow(uid int64, vid int64, action string) *favoriteFlow {
	return &favoriteFlow{uid: uid, vid: vid, action_type: action}
}

func (f *favoriteFlow) Do() (int64, error) {
	uid, err := f.favoriteAction()
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func (f *favoriteFlow) favoriteAction() (int64, error) {
	if f.action_type == "1" { // 如果操作类型是 点赞
		// video中的赞数量+1, is_favourite=true
		err := repository.NewVideoDao().AddVideoFavorite(f.vid) // 通过视频查询到video类
		if err != nil {
			return 0, nil
		}
		// 用户喜欢列表 video+1
		
	} else { // 如果操作类型是 取消点赞
		err := repository.NewVideoDao().RmVideoFavorite(f.vid) // 通过视频查询到video类
		if err != nil {
			return 0, nil
		}
	}
	return 0, nil
}
