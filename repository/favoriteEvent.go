package repository

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func SendFavoroteEventToKafka(uid int64, to_uid int64, action_type string) error {
	event := InteractiveEvent{
		Uid:         uid,
		To_uid:      to_uid,
		Action_type: "favorite+" + action_type,
	}

	eventBytes, _ := json.Marshal(event)

	return FavoriteWriter.WriteMessages(c, kafka.Message{
		Key:   nil,
		Value: eventBytes,
	})
}

// vid的点赞数量
func FavoriteUser(vid int64) error { // uid关注了to_uid
	// 关注列表
	favoritekey := fmt.Sprintf("favorite:%d", vid)
	exists, err := rdb4.Exists(c, favoritekey).Result()
	if err != nil {
		fmt.Println("Error checking existence...")
		return err
	}
	if exists != 1 { // 不存在这个键
		favoriteCount, err := NewFavoriteDao().CountFavoriteByVid(vid)
		if err != nil {
			return err
		}
		err = rdb4.Set(c, favoritekey, favoriteCount, 0).Err()
		if err != nil {
			return err
		}
	}
	// add code here
	_, err = rdb4.Incr(c, favoritekey).Result()
	if err != nil {
		return err
	}

	return nil
}

func UnFavoriteUser(vid int64) error {
	favoritekey := fmt.Sprintf("favorite:%d", vid)
	exists, err := rdb4.Exists(c, favoritekey).Result()
	if err != nil {
		fmt.Println("Error checking existence...")
		return err
	}
	if exists != 1 { // 不存在这个键
		favoriteCount, err := NewFavoriteDao().CountFavoriteByVid(vid)
		if err != nil {
			return err
		}
		err = rdb4.Set(c, favoritekey, favoriteCount, 0).Err()
		if err != nil {
			return err
		}
	}
	_, err = rdb4.Decr(c, favoritekey).Result()
	if err != nil {
		return err
	}

	return nil
}
