package repository

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func SendFollowEventToKafka(uid int64, to_uid int64, action_type string) error {
	event := InteractiveEvent{
		Uid:         uid,
		To_uid:      to_uid,
		Action_type: "follow+" + action_type,
	}

	eventBytes, _ := json.Marshal(event)

	return FollowWriter.WriteMessages(c, kafka.Message{
		Key:   nil,
		Value: eventBytes,
	})
}

// 关注用户:存在key，直接加入；不存在key，调出redis后加入
func FollowUser(uid int64, to_uid int64) error { // uid关注了to_uid
	// 关注列表
	followkey := fmt.Sprintf("following:%d", uid)
	exists, err := rdb3.Exists(c, followkey).Result()
	if err != nil {
		fmt.Println("Error checking existence...")
		return err
	}
	if exists != 1 { // 不存在这个键
		followlist, err := NewRalationDao().QueryFollow(uid)
		if err != nil {
			return err
		}
		if len(*followlist) > 0 {
			follows := make([]interface{}, len(*followlist))
			for i, relation := range *followlist {
				follows[i] = relation.ToUid
			}
			rdb3.SAdd(c, followkey, follows...)
		}
	}
	if err := StoreFollow(uid, to_uid); err != nil {
		return err
	}

	// 粉丝列表
	followerkey := fmt.Sprintf("follower:%d", to_uid)
	exists, err = rdb3.Exists(c, followerkey).Result()
	if err != nil {
		fmt.Println("Error checking existence...")
		return err
	}
	if exists != 1 {
		followlist, err := NewRalationDao().QueryFollower(to_uid)
		if err != nil {
			return err
		}
		if len(*followlist) > 0 {
			followers := make([]interface{}, len(*followlist))
			for i, relation := range *followlist {
				followers[i] = relation.ToUid
			}
			rdb3.SAdd(c, followerkey, followers...)
		}
	}
	if err := StoreFollower(to_uid, uid); err != nil {
		return err
	}
	return nil
}

func UnFollowUser(uid int64, to_uid int64) error {
	followkey := fmt.Sprintf("following:%d", uid)
	exists, err := rdb3.Exists(c, followkey).Result()
	if err != nil {
		fmt.Println("Error checking existence...")
		return err
	}
	if exists != 1 { // 不存在这个键
		followlist, err := NewRalationDao().QueryFollow(uid)
		if err != nil {
			return err
		}
		if len(*followlist) > 0 {
			follows := make([]interface{}, len(*followlist))
			for i, relation := range *followlist {
				follows[i] = relation.ToUid
			}
			rdb3.SRem(c, followkey, follows...)
		}
	}
	if err := RmFollow(uid, to_uid); err != nil {
		return err
	}

	// 粉丝列表
	followerkey := fmt.Sprintf("follower:%d", to_uid)
	exists, err = rdb3.Exists(c, followerkey).Result()
	if err != nil {
		fmt.Println("Error checking existence...")
		return err
	}
	if exists != 1 {
		followlist, err := NewRalationDao().QueryFollower(to_uid)
		if err != nil {
			return err
		}
		if len(*followlist) > 0 {
			followers := make([]interface{}, len(*followlist))
			for i, relation := range *followlist {
				followers[i] = relation.ToUid
			}
			rdb3.SAdd(c, followerkey, followers...)
		}
	}
	if err := RmFollower(to_uid, uid); err != nil {
		return err
	}
	return nil
}
