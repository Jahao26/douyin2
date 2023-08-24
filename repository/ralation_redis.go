package repository

import (
	"fmt"
	"strconv"
	"time"
)

// 存关注信息
func StoreFollow(uid int64, to_uid int64) error {
	FollowKey := fmt.Sprintf("following:%d", uid)
	Follow := strconv.FormatInt(to_uid, 10)
	fmt.Println("*********")
	fmt.Println(FollowKey, Follow)
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

// 存粉丝信息
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

	UserId := fmt.Sprintf("user:%d", uid)
	//pipe := rdb3.TxPipeline()
	//friendlist, err := pipe.SInter(c, FollowKey, FollowerKey).Result()
	//fmt.Println("redis friendlist:", friendlist)
	//pipe.Set(c, UserId, 0, 1*time.Minute)
	//_, err = pipe.Exec(c)
	friendlist, err := rdb3.SInter(c, FollowKey, FollowerKey).Result()
	rdb3.Set(c, UserId, 0, 1*time.Minute)
	if err != nil {
		return nil, err
	}
	fmt.Println("redis friendlist:", friendlist)

	return friendlist, err
}
