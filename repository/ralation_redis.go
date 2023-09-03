package repository

import (
	"fmt"
	"strconv"
	"time"
)

// 存uid关注信息
func StoreFollow(uid int64, to_uid int64) error {
	FollowKey := fmt.Sprintf("following:%d", uid)
	Follow := strconv.FormatInt(to_uid, 10)
	err := rdb3.SAdd(c, FollowKey, Follow).Err()
	return err
}

// 删除关注信息
func RmFollow(uid int64, to_uid int64) error {
	FollowKey := fmt.Sprintf("following:%d", uid)
	Follow := strconv.FormatInt(to_uid, 10)

	err := rdb3.SRem(c, FollowKey, Follow).Err()
	return err
}

// 获得关注列表
func GetFollow(uid int64) ([]string, error) {
	FollowKey := fmt.Sprintf("following:%d", uid)
	followlist, err := rdb3.SMembers(c, FollowKey).Result()
	if err != nil {
		return nil, err
	}
	return followlist, err
}

// 存uid粉丝信息
func StoreFollower(uid int64, to_uid int64) error {
	FollowerKey := fmt.Sprintf("follower:%d", uid)
	Follower := strconv.FormatInt(to_uid, 10)

	err := rdb3.SAdd(c, FollowerKey, Follower).Err()
	return err
}

// 删除粉丝信息
func RmFollower(uid int64, to_uid int64) error {
	FollowerKey := fmt.Sprintf("follower:%d", uid)
	Follower := strconv.FormatInt(to_uid, 10)

	err := rdb3.SRem(c, FollowerKey, Follower).Err()
	return err
}

// 获得粉丝列表
func GetFollower(uid int64) ([]string, error) {
	FollowerKey := fmt.Sprintf("follower:%d", uid)
	followerlist, err := rdb3.SMembers(c, FollowerKey).Result()
	if err != nil {
		return nil, err
	}
	return followerlist, err
}

// 获得朋友列表
func GetFriend(uid int64) ([]string, error) {
	FollowKey := fmt.Sprintf("following:%d", uid)
	FollowerKey := fmt.Sprintf("follower:%d", uid)

	exists, err := rdb3.Exists(c, FollowKey).Result()
	if err != nil {
		fmt.Println("Error checking friendList existence...")
		return nil, err
	}
	if exists != 1 { //不存在键
		// 先从数据库中取出朋友列表，再将他缓存在redis中
		// 方式1，如果我目前的联合查询有用，那么要在redis里存following和follower两个键值对要怎么存。
		// 方式2: 如果我直接用查询关注和查询粉丝的操作获得数据库中的两个列表，转存redis只能通过for进行遍历。
		fmt.Println("reload friend list!!!")
		followerlist, _ := NewRalationDao().QueryFollow(uid)
		for _, follower := range *followerlist {
			// 遍历关注者列表，挨个存入redis
			err = StoreFollow(uid, follower.ToUid)
		}
		followinglist, _ := NewRalationDao().QueryFollower(uid)
		for _, following := range *followinglist {
			// 遍历粉丝列表，存入redis
			// 由于mysql查粉丝列表时，当前用户uid作为to_uid进行查询，因此存redis时用的是Uid
			err = StoreFollower(uid, following.Uid)
		}
	}

	UserId := fmt.Sprintf("user:%d", uid)
	Usertime := fmt.Sprintf("usertime:%d", uid)
	//pipe := rdb3.TxPipeline()
	//friendlist, err := pipe.SInter(c, FollowKey, FollowerKey).Result()
	//fmt.Println("redis friendlist:", friendlist)
	//pipe.Set(c, UserId, 0, 1*time.Minute)
	//_, err = pipe.Exec(c)
	friendlist, err := rdb3.SInter(c, FollowKey, FollowerKey).Result()
	rdb3.Set(c, UserId, 0, 1*time.Minute)                   //设置聊天时间戳
	rdb3.Set(c, Usertime, time.Now().Unix(), 1*time.Minute) //设置聊天时间戳
	if err != nil {
		return nil, err
	}
	fmt.Println("redis friendlist:", friendlist)

	return friendlist, err
}
