package repository

import "sync"

type Ralation struct {
	Uid   int64 `gorm:"column:uid;not null" redis:"uid"`
	ToUid int64 `gorm:"column:to_uid;not null" redis:"to_uid"`
}

type RalationDAO struct {
}

var ralaionDao *RalationDAO
var ralationOnce sync.Once

func NewRalationDao() *RalationDAO {
	ralationOnce.Do(
		func() {
			ralaionDao = &RalationDAO{}
		})
	return ralaionDao
}

// 增加关注与被关注的关系
func (*RalationDAO) AddRalation(uid int64, to_uid int64) error {
	// 先添加对应关注信息到redis里，再存到mysql上。后续可通过消息队列优化，下同。
	err := StoreFollow(uid, to_uid)
	if err != nil {
		return err
	}
	err = StoreFollower(to_uid, uid)
	if err != nil {
		return err
	}

	newRala := Ralation{
		Uid:   uid,
		ToUid: to_uid,
	}
	if err := db.Create(&newRala).Error; err != nil {
		return err
	}
	return nil
}

// 删除关注与被关注的关系
func (*RalationDAO) RmRalation(uid int64, to_uid int64) error {
	err := RmFollow(uid, to_uid)
	if err != nil {
		return err
	}
	err = RmFollower(to_uid, uid)
	if err != nil {
		return err
	}

	var ralation Ralation
	if err := db.Where("uid=?", uid).Where("to_uid", to_uid).Delete(&ralation).Error; err != nil {
		return err
	}
	return nil
}

// 查看关注关系
func (*RalationDAO) QuaryRalation(uid int64, to_uid int64) (bool, error) {
	var rala Ralation
	err := db.Model(&Ralation{}).Where("uid=?", uid).Where("to_uid", to_uid).Find(&rala).Error
	if err != nil {
		return false, err
	}
	if rala.Uid == 0 || rala.ToUid == 0 {
		return false, err
	}
	return true, err
}

// 查询当前用户的关注列表
func (*RalationDAO) QuaryFollow(uid int64) (*[]Ralation, error) {
	var followList []Ralation
	err := db.Model(&Ralation{}).Where("uid", uid).Find(&followList).Error
	if err != nil {
		return nil, err
	}
	return &followList, err
}

// 查询当前用户的粉丝列表
func (*RalationDAO) QuaryFollower(uid int64) (*[]Ralation, error) {
	var followerList []Ralation
	err := db.Model(&Ralation{}).Where("to_uid", uid).Find(&followerList).Error
	if err != nil {
		return nil, err
	}
	return &followerList, err
}
