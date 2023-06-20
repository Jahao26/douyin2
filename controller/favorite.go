package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
// Demand: database(video-like)+1, 用户喜欢的视频列表+1
func FavoriteAction(c *gin.Context) {
	// token := c.Query("token")

	if _, exist := c.Get("uid"); exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
// Demand: 不同用户有不同的的喜欢视频的列表
func FavoriteList(c *gin.Context) {
	_, exist := c.Get("uid")
	if !exist {
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: []VideoResponse{},
	})
}
