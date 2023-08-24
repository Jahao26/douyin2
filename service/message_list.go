package service

import (
	"douyin/repository"
	"fmt"
)

var timestamp int64

type NewMassage struct {
	Id         int64  `json:"id,omitempty"`
	ToUid      int64  `json:"to_user_id,omitempty"`
	Uid        int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime int64  `json:"create_time,omitempty"`
}

type MassageList struct {
	Massages []*NewMassage
}

type QuaryMassageListFlow struct {
	uid         int64
	to_uid      int64
	massages    []*repository.Massage
	massagelist *MassageList
}

func QuaryMassageList(uid int64, to_uid int64) (*MassageList, error) {
	return newQuaryMassageListFlow(uid, to_uid).Do()
}

func newQuaryMassageListFlow(uid int64, to_uid int64) *QuaryMassageListFlow {
	return &QuaryMassageListFlow{uid: uid, to_uid: to_uid}
}

func (q *QuaryMassageListFlow) Do() (*MassageList, error) {
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.massagelist, nil
}

func (q *QuaryMassageListFlow) packData() error {

	err := repository.NewMassageDAO().QuaryMassage(q.uid, q.to_uid, &q.massages)
	if err != nil {
		return err
	}
	fmt.Println("out time:")

	newMassagelist := make([]*NewMassage, len(q.massages))
	for i := range q.massages {
		newmassage := NewMassage{
			Id:         q.massages[i].Id,
			Uid:        q.massages[i].Uid,
			ToUid:      q.massages[i].ToUid,
			Content:    q.massages[i].Content,
			CreateTime: q.massages[i].CreateTime,
			//CreateTime: q.massages[i].CreateTime.Format("2006-01-02 15:04:05"),
		}

		newMassagelist[i] = &newmassage
	}
	q.massagelist = &MassageList{Massages: newMassagelist}
	return nil
}
