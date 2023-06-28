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
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "failed to feed"},
		})
	}
	_, err := service.QueryUserById(userid.(int64))
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	// 获得推荐的视频列表，uid用于后续判断用户是否点赞
	uid := userid.(int64)
	videolist, err := service.QueryFeedVideoList(uid)

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videolist.Videos,
		NextTime:  time.Now().Unix(),
	})
}
