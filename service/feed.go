package service

import "douyin/repository"

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
