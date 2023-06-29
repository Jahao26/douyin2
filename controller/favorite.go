package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
// Demand: database(video-like)+1, 用户喜欢的视频列表+1
func FavoriteAction(c *gin.Context) {
	action_type := c.Query("action_type")
	video_id := c.Query("video_id")
	userid, exist := c.Get("uid")
	if exist {
		uid := userid.(int64)
		vid, _ := strconv.ParseInt(video_id, 10, 64)
		if err := service.FavoriteAction(uid, vid, action_type); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FavoriteList all users have same favorite video list
// Demand: 不同用户有不同的的喜欢视频的列表
func FavoriteList(c *gin.Context) {
	userid, exist := c.Get("uid")
	if !exist {
		return
	}
	uid := userid.(int64)
	_, err := service.QueryUserById(uid)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User not exist",
			},
			VideoList: []*service.VideoResponse{},
		})
	}
	//通过uid获得用户喜欢的视频列表
	videolist, err := service.FavoriteList(uid)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "favorite list failed",
			},
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videolist.Videos,
	})
}
