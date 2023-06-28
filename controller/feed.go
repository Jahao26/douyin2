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
			//{
			//	Id:     1,
			//	Author: user,
			//	//PlayUrl:  "https://vfx.mtime.cn/Video/2022/11/11/mp4/221111235609124131.mp4",
			//	//CoverUrl: "https://ww1.sinaimg.cn/mw690/007Sal7sly1henb0uz8z7j31jk25sb2a.jpg",
			//	PlayUrl:  "http://127.0.0.1:8080/static/1_1687833749_221111235609124131.mp4",
			//	CoverUrl: "http://127.0.0.1:8080/static/1_1687833749_221111235609124131.png",
			//
			//	FavoriteCount: 0,
			//	CommentCount:  0,
			//	IsFavorite:    false,
			//},
			//{
			//	Id:            2,
			//	Author:        user,
			//	PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
			//	CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
			//	FavoriteCount: 0,
			//	CommentCount:  0,
			//	IsFavorite:    false,
			//},
			//{
			//	Id:     3,
			//	Author: user,
			//	//PlayUrl: "http://127.0.0.1:8080/static/5_1687318088_VIDEO_20230621_112811122.mp4",
			//	//CoverUrl: "http://127.0.0.1:8080/static/5_1687318088_VIDEO_20230621_112811122.png",
			//	PlayUrl:       "./public/5_1687318088_VIDEO_20230621_112811122.mp4",
			//	CoverUrl:      "./public/5_1687318088_VIDEO_20230621_112811122.png",
			//	FavoriteCount: 0,
			//	CommentCount:  0,
			//	IsFavorite:    false,
			//},
			{
				Id:     4,
				Author: user,
				//PlayUrl:  "http://127.0.0.1:8080/static/5_1687318088_VIDEO_20230621_112811122.mp4",
				//CoverUrl: "http://127.0.0.1:8080/static/5_1687318088_VIDEO_20230621_112811122.png",
				PlayUrl:  "http://10.16.43.102:8080/static/5_1687318088_VIDEO_20230621_112811122.mp4",
				CoverUrl: "http://10.16.43.102:8080/static/5_1687318088_VIDEO_20230621_112811122.png",
				//PlayUrl:       "./public/5_1687318088_VIDEO_20230621_112811122.mp4",
				//CoverUrl:      "./public/5_1687318088_VIDEO_20230621_112811122.png",
				FavoriteCount: 0,
				CommentCount:  0,
				IsFavorite:    false,
			},
		},
		NextTime: time.Now().Unix(),
	})
}
