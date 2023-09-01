package repository

import (
	"fmt"
	"strconv"
)

// 存uid的喜欢视频号
func StoreVid(uid int64, vid int64) error {
	userKey := fmt.Sprintf("favorite:%d", uid)
	videoId := strconv.FormatInt(vid, 10)
	err := rdb3.SAdd(c, userKey, videoId).Err()
	return err
}

func RemVid(uid int64, vid int64) error {
	userKey := fmt.Sprintf("favorite:%d", uid)
	videoId := strconv.FormatInt(vid, 10)
	err := rdb3.SRem(c, userKey, videoId).Err()
	return err
}

// 获得Vid的点赞数量
func FavoriteCountByVid(vid int64) (int64, error) {
	favoritekey := fmt.Sprintf("favorite:%d", vid)
	count, err := rdb4.Get(c, favoritekey).Int64()
	if err != nil {
		return 0, err
	}
	return count, nil
}
