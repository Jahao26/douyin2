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
	IsFavorite    bool   `gorm:"column:is_favorite;" redis:"-"`
}

type VideoDAO struct {
}

var videoDao *VideoDAO
var videoOnce sync.Once

func NewVideoDao() *VideoDAO { // 单例设计模式：懒汉模式
	// 这个函数暴露唯一的接口对外使用
	userOnce.Do(
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

func (*VideoDAO) GetVideoByLimit(limit int, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList is nil")
	}
	if err := db.Limit(limit).Find(&videoList).Error; err != nil {
		return err
	}
	return nil
}

func (*VideoDAO) AddVideoFavorite(id int64) error { // 在增加视频喜欢的同时，将uid-vid-isfavorite更改
	if err := db.Model(&Video{}).Where("id=?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count+?", 1)); err != nil {
		panic(err)
	}
	if err := db.Model(&Video{}).Update("is_favorite", "true"); err != nil {
		panic(err)
	}
	return nil
}

func (*VideoDAO) RmVideoFavorite(id int64) error {
	if err := db.Model(&Video{}).Where("id=?", id).UpdateColumn("favorite_count", gorm.Expr("favorite_count-?", 1)); err != nil {
		panic(err)
	}
	if err := db.Model(&Video{}).Update("is_favorite", "false"); err != nil {
		panic(err)
	}
	return nil
}
