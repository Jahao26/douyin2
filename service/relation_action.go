package service

import (
	"douyin/repository"
	"fmt"
)

type relationFlow struct {
	uid         int64
	to_uid      int64
	action_type string
}

func RelationAction(uid int64, to_uid int64, action string) error {
	return newRelationFlow(uid, to_uid, action).Do()
}

func newRelationFlow(uid int64, to_uid int64, action string) *relationFlow {
	return &relationFlow{uid, to_uid, action}
}

func (r *relationFlow) Do() error {
	err := r.relationAction()
	if err != nil {
		return err
	}
	return nil
}

func (r *relationFlow) relationAction() error {
	if r.action_type == "1" { //如果是关注操作
		go func() {
			if err := repository.FollowUser(r.uid, r.to_uid); err != nil {
				fmt.Println("Error in Follow:", err)
			}
		}()

	} else { //取关操作
		go func() {
			if err := repository.UnFollowUser(r.uid, r.to_uid); err != nil {
				fmt.Println("Error in UnFollow:", err)
			}
		}()

	}
	go func() {
		if err := repository.SendFollowEventToKafka(r.uid, r.to_uid, r.action_type); err != nil {
			fmt.Println("Error in SendToKafka:", err)
		}
	}()

	return nil
}
