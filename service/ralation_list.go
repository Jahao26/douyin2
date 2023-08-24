package service

import (
	"douyin/repository"
	"strconv"
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
	if _, err := repository.GetUsr_redis(q.uid); err != nil {
		if _, err = repository.NewUserDao().QueryById(q.uid); err != nil {
			return err
		}
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

// 粉丝列表

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
	if _, err := repository.GetUsr_redis(q.uid); err != nil {
		if _, err = repository.NewUserDao().QueryById(q.uid); err != nil {
			return err
		}
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

// 朋友列表

func FriendList(uid int64) (*userList, error) {
	return NewFriendList(uid).friendDo()
}

func NewFriendList(uid int64) *RalationListFlow {
	return &RalationListFlow{uid: uid}
}

func (q *RalationListFlow) friendDo() (*userList, error) {
	if err := q.friendcheckNum(); err != nil {
		return nil, err
	}
	if err := q.friendpackData(); err != nil {
		return nil, err
	}
	return q.userlist, nil
}

func (q *RalationListFlow) friendcheckNum() error {
	//检查userId是否存在
	if _, err := repository.GetUsr_redis(q.uid); err != nil {
		if _, err = repository.NewUserDao().QueryById(q.uid); err != nil {
			return err
		}
	}
	return nil
}

func (q *RalationListFlow) friendpackData() error {
	// 通过uid获得朋友用户的信息列表
	// 1. 先判断redis中是否存在相关数据。(先不考虑关没关注)
	// 2. 若存在，则直接从redis中拉出来初始状态的朋友列表
	// 3. 若不存在，则从mysql中拉出来初始状态的朋友列表
	// 4. 对朋友列表补充关注信息。

	friList, err := repository.GetFriend(q.uid)

	if err != nil {
		return err
	}

	var newFriendList []*UserInfoPage
	for _, Friend := range friList { //获取关注列表里的信息用于返回
		Uid, _ := strconv.ParseInt(Friend, 10, 64)
		to_userinfo, _ := UserInfo(Uid)
		//查询关注与被关注的信息,改变显示的布尔值
		//change, _ := repository.NewRalationDao().QuaryRalation(q.uid, to_userinfo.Id)
		//to_userinfo.IsFollow = change

		newFriendList = append(newFriendList, to_userinfo)
	}
	q.userlist = &userList{newFriendList}
	return nil
}
