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

func InitRedis() error {
	rdb0 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       0,
	})
	rdb1 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       1,
	})
	rdb2 = redis.NewClient(&redis.Options{
		Addr:     "60.204.170.108:6379",
		Password: "",
		DB:       2,
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
