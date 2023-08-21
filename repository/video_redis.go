package repository

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// rdb1中存放视频信息，rdb2中存放用户历史观看key=user:id
var VIDEO_EXPIRE = 30 * time.Minute

// 加载user没看过的热点视频，放到redis里面
func Reload_redis(uid int64, limit int64, videoList *[]*Video) error {
	count, err := countRedis()
	if err != nil {
		fmt.Println("Error count Redis data: %v", err)
	}

	if count == 0 {
		if err := LoadUnvisitedVideosFromMysql(uid, limit, videoList); err != nil {
			return err
		}
	} else {
		UserUnvisitedVidList, _ := GetUnvisitedVideosFromRedis(uid) // 从redis获取user没有遍历过的视频id列表
		listLen := int64(len(UserUnvisitedVidList))
		if listLen <= limit { //不够limit个
			// 从SQL里面取出新热门视频（不在原来的redis缓存中）
			if err := NewVideoDao().GetNewVideos(limit-listLen, videoList); err != nil {
				return err
			}
		} else {
			// 取出没看过的所有视频
			/*Please add code here*/
			if err := GetVideosByVidFromRedis(UserUnvisitedVidList, videoList); err != nil {
				return err
			}
		}
	}

	return nil
}

// 存入单个视频
func StoreVideoIntoRedis(video *Video) error {
	videokey := fmt.Sprintf("video:%d", video.Id)
	pipe := rdb1.TxPipeline()

	videomap := map[string]interface{}{
		"id":             video.Id,
		"uid":            video.Uid,
		"playurl":        video.PlayUrl,
		"coverurl":       video.CoverUrl,
		"favorite_count": video.FavoriteCount,
		"comment_count":  video.CommentCount,
		"is_favorite":    video.IsFavorite,
	}

	pipe.HMSet(c, videokey, videomap)
	pipe.Expire(c, videokey, VIDEO_EXPIRE+time.Duration(rand.Float64()*EXPIRE_TIME_JITTER.Seconds())*time.Second)
	_, err := pipe.Exec(c)
	return err
}

// 通过遍历键获得数据数量的方法
func countRedis() (int64, error) {
	key := "video:*"
	keys, err := rdb1.Keys(c, key).Result()
	if err != nil {
		fmt.Println("Error in looking for keys")
		return 0, err
	}
	count := int64(len(keys))
	return count, nil
}

// map2Video函数，用来做redis map到video的转化
func mapToVideo(videomap map[string]string) (*Video, error) {
	video := &Video{}
	video.Id, _ = strconv.ParseInt(videomap["id"], 10, 64)
	video.Uid, _ = strconv.ParseInt(videomap["uid"], 10, 64)
	video.PlayUrl = videomap["playurl"]
	video.CoverUrl = videomap["coverurl"]
	video.FavoriteCount, _ = strconv.ParseInt(videomap["favorite_count"], 10, 64)
	video.CommentCount, _ = strconv.ParseInt(videomap["comment_count"], 10, 64)
	video.IsFavorite = videomap["is_favorite"] == "true"

	return video, nil
}

// 从mysql中取出limit个user没看过的视频存入redis
func LoadUnvisitedVideosFromMysql(uid int64, limit int64, videoList *[]*Video) error {
	if limit == 0 {
		return nil
	}
	// 取前limit个SQL数据，不重复，并且是实时的热点数据
	if err := NewVideoDao().GetHotVideos(uid, limit, videoList); err != nil {
		return err
	}
	// 将新获取的数据放到redis中
	for _, video := range *videoList {
		if err := StoreVideoIntoRedis(video); err != nil {
			fmt.Println("Error %v storing videos in video: %d", err, video.Id)
		}
	}

	return nil
}

// 从全部的视频set A中，获得set B中没有的集合
func GetUnvisitedVideosFromRedis(uid int64) ([]string, error) {
	setAkey := fmt.Sprintf("user:*")
	setBkey := fmt.Sprintf("user:%d", uid)
	diffMember, err := rdb2.SDiff(c, setAkey, setBkey).Result()
	if err != nil {
		return nil, err
	}
	return diffMember, nil
}

// 通过Vid获取视频，存到List中
func GetVideosByVidFromRedis(vidList []string, videoList *[]*Video) error {
	for _, vid := range vidList {
		videoKey := fmt.Sprintf("video:%s", vid)
		videoData, err := rdb1.HGetAll(c, videoKey).Result()
		if err != nil {
			return err
		}

		// 解析Video信息
		video, err := mapToVideo(videoData)
		*videoList = append(*videoList, video)
	}

	return nil
}

func StoreVidIntoRedis(uid int64, vidList []int64) error {
	key := fmt.Sprintf("user:%d", uid)
	vidStrings := make([]string, len(vidList))
	for i, vid := range vidList {
		vidStrings[i] = strconv.FormatInt(vid, 10)
	}
	err := rdb2.SAdd(c, key, vidStrings).Err()
	return err
}
