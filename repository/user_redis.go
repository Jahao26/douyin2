package repository

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// 设置上下文
var c = context.Background()

// 过期时间
var (
	USER_INFO_EXPIRE   = 20 * time.Minute
	EXPIRE_TIME_JITTER = 5 * time.Minute
)

func AddUsrInfo_redis(user *User) error {
	userKey := fmt.Sprintf("user:%d", user.Id)
	pipe := rdb0.TxPipeline()

	userinfoMap := map[string]interface{}{
		"id":             user.Id,
		"name":           user.Name,
		"follow_count":   user.FollowCount,
		"follower_count": user.FollowerCount,
		"is_follow":      user.IsFollow,
		"favorite_count": user.FavoriteCount,
		"work_count":     user.WorkCount,
	}
	// 向pipeline中写入数据
	pipe.HMSet(c, userKey, userinfoMap)

	// 设置过期时间，此处EXPIRE_TIME_JITTER保证了所有userinfo键不同时过期
	pipe.Expire(c, userKey, USER_INFO_EXPIRE+time.Duration(rand.Float64()*EXPIRE_TIME_JITTER.Seconds())*time.Second)

	// 执行pipeline事务
	_, err := pipe.Exec(c)
	return err
}

func GetUsr_redis(uid int64) (*User, error) {
	userKey := fmt.Sprintf("user:%d", uid)
	userinfoMap, err := rdb0.HGetAll(c, userKey).Result()
	if err != nil {
		return nil, err
	}

	if len(userinfoMap) == 0 {
		return nil, fmt.Errorf("User with ID %d not found in Redis", uid)
	}
	user := &User{}
	user.Id = uid
	user.Name = userinfoMap["name"]
	user.FollowCount, _ = strconv.ParseInt(userinfoMap["follow_count"], 10, 64)
	user.FollowerCount, _ = strconv.ParseInt(userinfoMap["follower_count"], 10, 64)
	user.IsFollow = userinfoMap["is_follow"] == "true"
	user.FavoriteCount, _ = strconv.ParseInt(userinfoMap["favorite_count"], 10, 64)
	user.WorkCount, _ = strconv.ParseInt(userinfoMap["work_count"], 10, 64)

	return user, nil
}
