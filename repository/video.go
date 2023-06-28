package repository

import (
	"errors"
	"sync"
)

type Video struct {
	Id            int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" redis:"id"`
	Uid           int64  `gorm:"column:uid;not null" redis:"uid"`
	PlayUrl       string `gorm:"column:playurl;type:varchar(255);not null" redis:"playurl"`
	CoverUrl      string `gorm:"column:coverurl;type:varchar(255);not null" redis:"coverurl"`
	FavoriteCount int64  `gorm:"column:favorite_count;not null" redis:"favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count;not null" redis:"comment_count"`
	IsFavorite    bool   `gorm:"column:is_favourite;" redis:"-"`
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

func (*VideoDAO) GetVideo(uid int64, videoList *[]*Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId videoList 空指针")
	}
	return db.Where("uid=?", uid).
		Find(&videoList).Error
}
