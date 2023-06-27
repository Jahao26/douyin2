package controller

import (
	"douyin/service"
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
	userid, exist := c.Get("uid")
	if !exist {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "failed to feed"},
		})
	}
	user, err := service.QueryUserById(userid.(int64))
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	// 统一传入空列表，没有视频
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{StatusCode: 0},
		VideoList: []VideoResponse{
		{
			Id:            1,
			Author:        user,
			PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
			CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
		},
		},
		NextTime: time.Now().Unix(),
	})
}
