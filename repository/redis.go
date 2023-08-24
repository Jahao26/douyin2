package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb0 *redis.Client // 操作用户数据
var rdb1 *redis.Client // 存储视频数据
var rdb2 *redis.Client //
var rdb3 *redis.Client // 存储关注、粉丝、好友列表

func InitRedis() error {
	// rdb0存放用户基本信息
	rdb0 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       0,
	})
	// rdb1存放视频缓存信息
	rdb1 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       1,
	})
	// rdb2存放用户历史观看信息
	rdb2 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       2,
	})
	// rdb3存放用户关注、粉丝、好友信息
	rdb3 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       3,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb0.Ping(ctx).Result()
	if err != nil {
		return errors.New(err.Error())
	}
	fmt.Println("****************************************")
	fmt.Println(pong)
	fmt.Println("****************************************")
	return nil
}
