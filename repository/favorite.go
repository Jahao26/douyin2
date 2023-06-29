package repository

import (
	"sync"
)

type Favorite struct {
	Uid        int64 `gorm:"column:uid;not null" redis:"uid"`
	VId        int64 `gorm:"column:vid;primary_key;AUTO_INCREMENT" redis:"vid"`
	IsFavorite bool  `gorm:"column:is_favorite;not null" redis:"-"`
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
		Uid:        uid,
		VId:        vid,
		IsFavorite: true,
	}
	if err := db.Create(&newFav).Error; err != nil {
		return err
	}
	return nil
}

func (*FavoriteDAO) RmFavorite(uid int64, vid int64) error {
	var favorite Favorite
	err := db.Where("uid=?", uid).Where("vid=?", vid).Find(&favorite).Error
	if err != nil {
		return err
	}
	if err = db.Delete(&favorite).Error; err != nil {
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
