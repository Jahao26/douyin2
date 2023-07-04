package service

import "douyin/repository"

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
		//根据当前用户id和被关注用户id，增加关系列表
		err := repository.NewRalationDao().AddRalation(r.uid, r.to_uid)
		if err != nil {
			return err
		}
		err = repository.NewUserDao().AddFollow(r.uid) //当前用户关注+1
		if err != nil {
			return err
		}
		err = repository.NewUserDao().AddFollower(r.to_uid) //被关注用户粉丝+1
		if err != nil {
			return err
		}
	} else { //取关操作
		// 根据当前用户id和被关注用户id，删除关系列表
		err := repository.NewRalationDao().RmRalation(r.uid, r.to_uid)
		if err != nil {
			return err
		}
		err = repository.NewUserDao().RmFollow(r.uid) //当前用户关注-1
		if err != nil {
			return err
		}
		err = repository.NewUserDao().RmFollower(r.to_uid) //被关注用户粉丝-1
		if err != nil {
			return err
		}
	}
	return nil
}
