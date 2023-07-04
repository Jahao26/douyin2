package service

import (
	"douyin/repository"
	"errors"
)

// userList 视频集合
type userList struct {
	Users []*UserInfoPage `json:"user_list,omitempty"`
}

type RalationListFlow struct {
	uid      int64
	users    []*UserInfoPage
	userlist *userList
}

func FollowList(uid int64) (*userList, error) {
	return NewFollowList(uid).followDo()
}

func NewFollowList(uid int64) *RalationListFlow {
	return &RalationListFlow{uid: uid}
}

func (q *RalationListFlow) followDo() (*userList, error) {
	if err := q.followChecknum(); err != nil {
		return nil, err
	}
	if err := q.followPackdata(); err != nil {
		return nil, err
	}
	return q.userlist, nil
}

func (q *RalationListFlow) followChecknum() error {
	//检查userId是否存在
	_, err := repository.NewUserDao().QueryById(q.uid)
	if err != nil {
		return errors.New("用户不存在")
	}
	return nil
}

func (q *RalationListFlow) followPackdata() error {
	// 通过uid获得关注用户的信息列表
	folList, err := repository.NewRalationDao().QuaryFollow(q.uid)
	if err != nil {
		return err
	}
	var newFollowList []*UserInfoPage
	for _, follow := range *folList { //获取关注列表里的信息用于返回
		to_userinfo, _ := UserInfo(follow.ToUid) //通过列表中被关注者的id，获得被关注者的用户信息
		to_userinfo.IsFollow = true              //在关注列表里的，都显示已关注
		newFollowList = append(newFollowList, to_userinfo)
	}
	q.userlist = &userList{newFollowList}
	return nil
}

func FollowerList(uid int64) (*userList, error) {
	return NewFollowerList(uid).Do()
}

func NewFollowerList(uid int64) *RalationListFlow {
	return &RalationListFlow{uid: uid}
}

func (q *RalationListFlow) Do() (*userList, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.userlist, nil
}

func (q *RalationListFlow) checkNum() error {
	//检查userId是否存在
	_, err := repository.NewUserDao().QueryById(q.uid)
	if err != nil {
		return errors.New("用户不存在")
	}
	return nil
}

func (q *RalationListFlow) packData() error {
	// 通过uid获得粉丝用户的信息列表
	folList, err := repository.NewRalationDao().QuaryFollower(q.uid)
	if err != nil {
		return err
	}
	var newFollowerList []*UserInfoPage
	for _, follower := range *folList { //获取关注列表里的信息用于返回
		to_userinfo, _ := UserInfo(follower.Uid)
		//查询关注与被关注的信息,改变显示的布尔值
		change, _ := repository.NewRalationDao().QuaryRalation(q.uid, to_userinfo.Id)
		to_userinfo.IsFollow = change

		newFollowerList = append(newFollowerList, to_userinfo)
	}
	q.userlist = &userList{newFollowerList}
	return nil
}
