package service

import (
	"douyin/repository"
	"time"
)

type massageFlow struct {
	uid         int64
	to_uid      int64
	action_type string
	content     string
}

func MassageAction(uid int64, to_uid int64, action_type string, content string) error {
	return newMassageFlow(uid, to_uid, action_type, content).Do()
}

func newMassageFlow(uid int64, to_uid int64, action_type string, content string) *massageFlow {
	return &massageFlow{uid, to_uid, action_type, content}
}

func (m *massageFlow) Do() error {
	err := m.massageAction()
	if err != nil {
		return err
	}
	return nil
}

func (m *massageFlow) massageAction() error {
	if m.action_type == "1" {
		// 当发送消息时，通过数据库操作将消息信息保存在数据库中
		newMassage := &repository.Massage{
			Uid:        m.uid,
			ToUid:      m.to_uid,
			Content:    m.content,
			CreateTime: time.Now().Unix(),
		}
		err := repository.NewMassageDAO().AddMassage(newMassage)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
