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

	//创建一个新的视频列表，长度和查询到的视频列表一致，用来返回给前端
	newvideolist := make([]*VideoResponse, len(q.videos))

	for i := range q.videos {
		newResponse := VideoResponse{
			Id:            q.videos[i].Id,
			Author:        user,
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

// feed列表

// MaxVideoNum 每次最多返回的视频数量
const (
	MaxVideoNum = 30
)

func QueryFeedVideoList(uid int64) (*List, error) {
	return NewQueryFeedVideoListFlow(uid).Do()
}

func NewQueryFeedVideoListFlow(uid int64) *QueryFeedVideoListFlow {
	return &QueryFeedVideoListFlow{uid: uid}
}

type QueryFeedVideoListFlow struct {
	uid    int64
	videos []*repository.Video

	videoList *List
}

func (q *QueryFeedVideoListFlow) Do() (*List, error) {
	//所有传入的参数不填也应该给他正常处理
	q.checkNum()

	if err := q.prepareData(); err != nil {
		return nil, err
	}

	return q.videoList, nil
}

func (q *QueryFeedVideoListFlow) checkNum() {
	//上层通过把userId置零，表示userId不存在或不需要
	if q.uid > 0 {
		//这里说明userId是有效的，可以定制性的做一些登录用户的专属视频推荐
	}
}

func (q *QueryFeedVideoListFlow) prepareData() error {
	err := repository.NewVideoDao().GetVideoByLimit(MaxVideoNum, &q.videos)
	if err != nil {
		return err
	}
	//如果用户为登录状态，则更新该视频是否被该用户点赞的状态
	//此处预留一个部分用于更新视频状态

	user, err := repository.NewUserDao().QueryById(q.uid)

	//创建一个新的视频列表，长度和查询到的视频列表一致，用来返回给前端
	newvideolist := make([]*VideoResponse, len(q.videos))

	for i := range q.videos {
		newResponse := VideoResponse{
			Id:            q.videos[i].Id,
			Author:        user,
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
