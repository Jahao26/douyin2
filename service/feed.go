package service

import (
	"douyin/repository"
	"fmt"
)

// feed列表
type FeedVideoListFlow struct {
	uid       int64
	videos    []*repository.Video
	videoList *List
}

// MaxVideoNum 每次最多返回的视频数量
const (
	MaxVideoNum = 30
)

func FeedVideoList(uid int64) (*List, error) {
	return NewFeedVideoListFlow(uid).Do()
}

func NewFeedVideoListFlow(uid int64) *FeedVideoListFlow {
	return &FeedVideoListFlow{uid: uid}
}

func (q *FeedVideoListFlow) Do() (*List, error) {
	//所有传入的参数不填也应该给他正常处理
	q.checkNum()

	if err := q.prepareData(); err != nil {
		return nil, err
	}

	return q.videoList, nil
}

func (q *FeedVideoListFlow) checkNum() {
	//上层通过把userId置零，表示userId不存在或不需要
	if q.uid > 0 {
		//这里说明userId是有效的，可以定制性的做一些登录用户的专属视频推荐
	}
}

func (q *FeedVideoListFlow) prepareData() error {
	// 重新加载热点数据
	if err := repository.Reload_redis(q.uid, MaxVideoNum, &q.videos); err != nil {
		return err
	}

	err := repository.NewVideoDao().GetVideoByLimit(MaxVideoNum, &q.videos)
	if err != nil {
		return err
	}

	//创建一个新的视频列表，长度和查询到的视频列表一致，用来返回给前端
	newvideolist := make([]*VideoResponse, len(q.videos))
	var vidList []int64
	var is_fav bool

	for i := range q.videos {
		is_fav, _ = repository.NewFavoriteDao().QueryUidVid(q.uid, q.videos[i].Id)
		// 通过保存的视频的UID找到作者信息，返回给feed列表
		author, _ := UserInfo(q.videos[i].Uid)
		is_rala, _ := repository.NewRalationDao().QuaryRalation(q.uid, author.Id)

		author.IsFollow = is_rala

		newResponse := VideoResponse{
			Id:            q.videos[i].Id,
			Author:        author,
			CommentCount:  q.videos[i].CommentCount,
			PlayUrl:       q.videos[i].PlayUrl,
			CoverUrl:      q.videos[i].CoverUrl,
			FavoriteCount: q.videos[i].FavoriteCount,
			IsFavorite:    is_fav,
		}
		newvideolist[i] = &newResponse
		vidList = append(vidList, q.videos[i].Id)
	}
	q.videoList = &List{Videos: newvideolist}
	//每次加载完视频都会默认看完了这些视频

	if err := repository.StoreVidIntoRedis(q.uid, vidList); err != nil {
		fmt.Println("Error in Add vid into Redis")
	}
	return nil
}
