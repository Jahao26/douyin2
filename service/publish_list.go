package service

import (
	"douyin/repository"
)

type VideoResponse struct {
	Id            int64         `json:"id,omitempty"`
	Author        *UserInfoPage `json:"author,omitempty"`
	PlayUrl       string        `json:"play_url,omitempty"`
	CoverUrl      string        `json:"cover_url,omitempty"`
	FavoriteCount int64         `json:"favorite_count,omitempty"`
	CommentCount  int64         `json:"comment_count,omitempty"`
	IsFavorite    bool          `json:"is_favorite,omitempty"`
}

// List 视频集合
type List struct {
	Videos []*VideoResponse `json:"video_list,omitempty"`
}

type QueryVideoListByUserIdFlow struct {
	uid    int64
	videos []*repository.Video

	videoList *List
}

func QuaryVideolistByUid(uid int64) (*List, error) {
	return NewQueryVideoListByUserIdFlow(uid).Do()
}

func NewQueryVideoListByUserIdFlow(uid int64) *QueryVideoListByUserIdFlow {
	return &QueryVideoListByUserIdFlow{uid: uid}
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
	if _, err := repository.GetUsr_redis(q.uid); err != nil {
		if _, err = repository.NewUserDao().QueryById(q.uid); err != nil {
			return err
		}
	}
	return nil
}

// 注意：Video由于在数据库中没有存储作者信息，所以需要手动填充
// 拼接将放在controller完成
func (q *QueryVideoListByUserIdFlow) packData() error {
	//通过uid获得视频列表，并放入视频流的videos
	err := repository.NewVideoDao().GetVideoByUid(q.uid, &q.videos)
	if err != nil {
		return err
	}

	userInfo, err := UserInfo(q.uid) // need Author info

	//创建一个新的视频列表，长度和查询到的视频列表一致，用来返回给前端
	newvideolist := make([]*VideoResponse, len(q.videos))

	for i := range q.videos {
		newResponse := VideoResponse{
			Id:            q.videos[i].Id,
			Author:        userInfo,
			CommentCount:  q.videos[i].CommentCount,
			PlayUrl:       q.videos[i].PlayUrl,
			CoverUrl:      q.videos[i].CoverUrl,
			FavoriteCount: q.videos[i].FavoriteCount,
			IsFavorite:    q.videos[i].IsFavorite,
		}
		newvideolist[i] = &newResponse
	}
	q.videoList = &List{Videos: newvideolist}
	return nil
}
