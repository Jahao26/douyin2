package service

import (
	"douyin/repository"
	"errors"
)

type FavoriteListFlow struct {
	uid       int64
	videos    []*repository.Video
	videoList *List
}

func FavoriteList(uid int64) (*List, error) {
	return NewFavoriteList(uid).Do()
}

func NewFavoriteList(uid int64) *FavoriteListFlow {
	return &FavoriteListFlow{uid: uid}
}

func (q *FavoriteListFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *FavoriteListFlow) checkNum() error {
	//检查userId是否存在
	_, err := repository.NewUserDao().QueryById(q.uid)
	if err != nil {
		return errors.New("用户不存在")
	}
	return nil
}

func (q *FavoriteListFlow) packData() error {
	//通过uid获得喜欢的视频列表
	favList, err := repository.NewFavoriteDao().GetFavoriteByUid(q.uid)
	if err != nil {
		return err
	}
	var newResponseList []*VideoResponse

	for _, fav := range *favList { // 获取喜欢视频列表的视频及作者信息
		video, _ := repository.NewVideoDao().GetVideoByVid(fav.VId)
		author, _ := repository.NewUserDao().QueryById(fav.Uid)

		newResponse := VideoResponse{
			Id:            video.Id,
			Author:        author,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
		}
		newResponseList = append(newResponseList, &newResponse)

	}
	q.videoList = &List{Videos: newResponseList}
	return nil
}
