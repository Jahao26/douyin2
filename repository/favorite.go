package repository

import (
	"sync"
)

type Favorite struct {
	Uid int64 `gorm:"column:uid;not null" redis:"uid"`
	VId int64 `gorm:"column:vid;not null" redis:"vid"`
}

type FavoriteDAO struct {
}

var favoriteDao *FavoriteDAO
var favoriteOnce sync.Once

func NewFavoriteDao() *FavoriteDAO {
	favoriteOnce.Do(
		func() {
			favoriteDao = &FavoriteDAO{}
		})
	return favoriteDao
}

func (*FavoriteDAO) AddFavorite(uid int64, vid int64) error {
	newFav := Favorite{
		Uid: uid,
		VId: vid,
	}
	if err := db.Create(&newFav).Error; err != nil {
		return err
	}
	return nil
}

func (*FavoriteDAO) RmFavorite(uid int64, vid int64) error {
	var favorite Favorite
	err := db.Where("uid=?", uid).Where("vid=?", vid).Delete(&favorite).Error
	if err != nil {
		return err
	}
	return nil
}

func (*FavoriteDAO) QueryUidVid(uid int64, vid int64) (bool, error) {
	var fav Favorite
	err := db.Model(&Favorite{}).Where("uid=?", uid).Where("vid=?", vid).Find(&fav).Error
	if err != nil {
		return false, err
	}
	if fav.Uid == 0 {
		return false, nil
	}
	return true, nil
}

// 通过UID查询喜欢的视频列表
func (*FavoriteDAO) GetFavoriteByUid(uid int64) (*[]Favorite, error) {
	// 通过uid查到所有的视频
	var favList []Favorite
	err := db.Model(&Favorite{}).Where("uid=?", uid).Find(&favList).Error
	if err != nil {
		return nil, err
	}
	return &favList, nil
}
