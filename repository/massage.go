package repository

import (
	"sync"
)

type Massage struct {
	Id         int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" redis:"id"`            //消息id
	Uid        int64  `gorm:"column:uid;not null" redis:"uid"`                            //发送者id
	ToUid      int64  `gorm:"column:to_uid;not null" redis:"to_uid"`                      //接收者id
	Content    string `gorm:"column:content;type:varchar(1000);not null" redis:"content"` //消息内容
	CreateTime int64  `gorm:"column:create_time;not null" redis:"create_time"`            //创建时间
}

type MassageDAO struct {
}

var massageDao *MassageDAO
var massageOnce sync.Once

func NewMassageDAO() *MassageDAO {
	massageOnce.Do(
		func() {
			massageDao = &MassageDAO{}
		})
	return massageDao
}

// 增加消息记录
func (m *MassageDAO) AddMassage(massage *Massage) error {
	if err := db.Create(&massage).Error; err != nil {
		return err
	}
	return nil
}

// 通过uid和to_uid获取消息记录
func (m *MassageDAO) QuaryMassage(uid int64, to_uid int64, timestamp int64, massageList *[]*Massage) (int64, error) {

	err := db.Model(&Massage{}).Where("create_time>?", timestamp).Where("(uid=? AND to_uid =?) OR (uid=? AND to_uid =?)", uid, to_uid, to_uid, uid).Find(&massageList).Error
	if err != nil {
		return 0, err
	}
	err = db.Model(&Massage{}).Select("MAX(create_time)").Scan(&timestamp).Error

	return timestamp, nil
}
