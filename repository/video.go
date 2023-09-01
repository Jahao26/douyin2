package repository

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type Video struct {
	Id            int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" redis:"id"`
	Uid           int64  `gorm:"column:uid;not null" redis:"uid"`
	PlayUrl       string `gorm:"column:playurl;type:varchar(255);not null" redis:"playurl"`
	CoverUrl      string `gorm:"column:coverurl;type:varchar(255);not null" redis:"coverurl"`
	FavoriteCount int64  `gorm:"column:favorite_count;not null" redis:"favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count;not null" redis:"comment_count"`
	IsFavorite    bool   `gorm:"column:is_favorite;not null" redis:"is_favorite"`
}

type VideoDAO struct {
}

var videoDao *VideoDAO
var videoOnce sync.Once

func NewVideoDao() *VideoDAO { // 单例设计模式：懒汉模式
	// 这个函数暴露唯一的接口对外使用
	videoOnce.Do(
		func() {
			videoDao = &VideoDAO{}
		})
	return videoDao
}

func (*VideoDAO) AddVideo(video *Video) error {
	if err := db.Create(&video).Error; err != nil {
		return err
	}
	return nil
}

func (*VideoDAO) GetVideoByUid(uid int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList is nil")
	}
	return db.Where("uid=?", uid).
		Find(&videoList).Error
}

func (*VideoDAO) GetVideoByVid(vid int64) (*Video, error) {
	var video Video
	err := db.Model(Video{}).Where("id=?", vid).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (*VideoDAO) GetVideoByLimit(limit int, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList is nil")
	}
	if err := db.Order("RAND()").Limit(limit).Find(&videoList).Error; err != nil {
		return err
	}
	return nil
}

func (*VideoDAO) AddVideoFavorite(vid int64) error { // 在增加视频喜欢的同时，将uid-vid-isfavorite更改
	db.Model(&Video{}).Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count+?", 1))
	//db.Model(&Video{}).Where("id=?", vid).Where("uid=?", uid).Update("is_favorite", true)
	return nil
}

func (*VideoDAO) RmVideoFavorite(vid int64) error {
	db.Model(&Video{}).Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count-?", 1))
	//db.Model(&Video{}).Where("id=?", vid).Where("uid=?", uid).Update("is_favorite", false)
	return nil
}

func (*VideoDAO) AddVideoComment(vid int64) error {
	db.Model(&Video{}).Where("id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count+?", 1))
	return nil
}

func (*VideoDAO) RmVideoComment(vid int64) error {
	db.Model(&Video{}).Where("id=?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count-?", 1))
	return nil
}

func (*VideoDAO) GetHotVideos(uid int64, limit int64, videoList *[]*Video) error { // 获取前limit个不重复的视频数据,存放在videolist中
	var visitedVideos []string
	key := fmt.Sprintf("user:%d", uid)
	if err := rdb2.SMembers(c, key).ScanSlice(&visitedVideos); err != nil {
		return err
	}

	// 构建查询
	query := db.Model(&Video{}).Where("id NOT IN (?)", visitedVideos).
		Order("favorite_count + comment_count DESC").
		Limit(int(limit)).Find(videoList)

	if query.Error != nil {
		return query.Error
	}

	return nil
}

func (*VideoDAO) GetNewVideos(limit int64, videoList *[]*Video) error { // 获取不在redis中的视频
	key := fmt.Sprintf("user:*")
	vidList, err := rdb2.SMembers(c, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	var excludedIDs []int64
	for _, idStr := range vidList {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		excludedIDs = append(excludedIDs, id)
	}

	query := db.Not(excludedIDs).Limit(int(limit)).Find(&videoList)
	if query.Error != nil {
		return query.Error
	}

	return nil
}
