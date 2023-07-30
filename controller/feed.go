package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []*service.VideoResponse `json:"video_list,omitempty"`
	NextTime  int64                    `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	userid, exist := c.Get("uid")
	if !exist {
		println("****user not exist!****")
		videolist, _ := service.FeedVideoList(int64(0))
		println(videolist.Videos[1].PlayUrl)
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videolist.Videos,
			NextTime:  time.Now().Unix(),
		})
		//c.JSON(http.StatusOK, FeedResponse{
		//	Response: Response{StatusCode: 1, StatusMsg: "failed to feed"},
		//})
	} else {
		_, err := service.UserInfo(userid.(int64))
		if err != nil {
			println("****lalalal")
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		// 获得推荐的视频列表，uid用于后续判断用户是否点赞
		uid := userid.(int64)
		videolist, err := service.FeedVideoList(uid)

		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videolist.Videos,
			NextTime:  time.Now().Unix(),
		})
	}
}
