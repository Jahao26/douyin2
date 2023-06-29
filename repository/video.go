package repository

import (
	"errors"
	"gorm.io/gorm"
	"sync"
)

type Video struct {
	Id            int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" redis:"id"`
	Uid           int64  `gorm:"column:uid;not null" redis:"uid"`
	PlayUrl       string `gorm:"column:playurl;type:varchar(255);not null" redis:"playurl"`
	CoverUrl      string `gorm:"column:coverurl;type:varchar(255);not null" redis:"coverurl"`
	FavoriteCount int64  `gorm:"column:favorite_count;not null" redis:"favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count;not null" redis:"comment_count"`
	IsFavorite    bool   `gorm:"column:is_favorite;not null" redis:"-"`
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
	if err := db.Limit(limit).Find(&videoList).Error; err != nil {
		return err
	}
	return nil
}

func (*VideoDAO) AddVideoFavorite(uid int64, vid int64) error { // 在增加视频喜欢的同时，将uid-vid-isfavorite更改
	db.Model(&Video{}).Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count+?", 1))
	//db.Model(&Video{}).Where("id=?", vid).Where("uid=?", uid).Update("is_favorite", true)
	return nil
}

func (*VideoDAO) RmVideoFavorite(uid int64, vid int64) error {
	db.Model(&Video{}).Where("id=?", vid).UpdateColumn("favorite_count", gorm.Expr("favorite_count-?", 1))
	//db.Model(&Video{}).Where("id=?", vid).Where("uid=?", uid).Update("is_favorite", false)
	return nil
}
