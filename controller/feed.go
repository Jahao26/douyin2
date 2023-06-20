package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []VideoResponse `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//userid, exist := c.Get("uid")
	//if !exist {
	//	c.JSON(http.StatusOK, FeedResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "failed to feed"},
	//	})
	//}
	// 统一传入空列表，没有视频
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: []VideoResponse{},
		NextTime:  time.Now().Unix(),
	})
}
