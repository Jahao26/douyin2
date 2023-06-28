package service

import (
	"douyin/repository"
	"errors"
)

type VideoResponse struct {
	Id            int64            `json:"id,omitempty"`
	Author        *repository.User `json:"author"`
	PlayUrl       string           `json:"play_url,omitempty"`
	CoverUrl      string           `json:"cover_url,omitempty"`
	FavoriteCount int64            `json:"favorite_count,omitempty"`
	CommentCount  int64            `json:"comment_count,omitempty"`
	IsFavorite    bool             `json:"is_favorite,omitempty"`
}

// List 视频集合
type List struct {
	//Videos []*repository.Video `json:"video_list,omitempty"`
	Videos []*VideoResponse `json:"video_list,omitempty"`
}

func QuaryVideolistByUid(uid int64) (*List, error) {
	return NewQueryVideoListByUserIdFlow(uid).Do()
}

func NewQueryVideoListByUserIdFlow(uid int64) *QueryVideoListByUserIdFlow {
	return &QueryVideoListByUserIdFlow{uid: uid}
}

type QueryVideoListByUserIdFlow struct {
	uid    int64
	videos []*repository.Video

	videoList *List
}

func (q *QueryVideoListByUserIdFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryVideoListByUserIdFlow) checkNum() error {
	//检查userId是否存在
	_, err := repository.NewUserDao().QueryById(q.uid)
	if err != nil {
		return errors.New("用户不存在")
	}
	return nil
}

// 注意：Video由于在数据库中没有存储作者信息，所以需要手动填充
// 此处由于Video格式不同，暂时只单纯返回列表
// 拼接将放在controller完成
func (q *QueryVideoListByUserIdFlow) packData() error {
	//通过uid获得视频列表，并放入视频流的videos
	err := repository.NewVideoDao().GetVideoByUid(q.uid, &q.videos)
	if err != nil {
		return err
	}
	user, err := repository.NewUserDao().QueryById(q.uid)
	var newvideolist = []*VideoResponse{}

	for i := range q.videos {
		newvideolist[i].Id = q.videos[i].Id
		newvideolist[i].Author = user
		newvideolist[i].CommentCount = q.videos[i].CommentCount
		newvideolist[i].PlayUrl = q.videos[i].PlayUrl
		newvideolist[i].CoverUrl = q.videos[i].CoverUrl
		newvideolist[i].FavoriteCount = q.videos[i].FavoriteCount
		newvideolist[i].IsFavorite = q.videos[i].IsFavorite
	}
	q.videoList = &List{Videos: newvideolist}
	return nil
}
