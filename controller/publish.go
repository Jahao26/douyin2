package controller

import (
	"douyin/repository"
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type VideoResponse struct {
	Id            int64            `json:"id,omitempty"`
	Author        *repository.User `json:"author"`
	PlayUrl       string           `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string           `json:"cover_url,omitempty"`
	FavoriteCount int64            `json:"favorite_count,omitempty"`
	CommentCount  int64            `json:"comment_count,omitempty"`
	IsFavorite    bool             `json:"is_favorite,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []VideoResponse `json:"video_list"`
}

// Publish check token then save upload file to public directory
// Action无法获取到token，post401
func Publish(c *gin.Context) {
	userid, exist := c.Get("uid") // 此处已经通过解析 token获得uid
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Uid get error",
		})
	}

	uid := userid.(int64)
	// FormFile源码：通过Request.ParseMultipartForm对上传文件参数进行解析，然后调用Request.FormFile获取文件头FileHeader
	// 获得上传的文件
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return

	}
	// 以用户id+文件名作为存储名，上传到本地
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%d_%s", uid, time.Now().Unix(), filename)
	// saveFile := filepath.Join("./public/", finalName)
	saveFile := "./public/" + finalName
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 定义封面名字
	strArray := strings.Split(finalName, ".")
	imageName := strArray[0]
	// 获取封面
	coverPath, err := service.GetCoverimage(saveFile, imageName)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	// 将视频信息存储到数据库
	// 将视频信息存储到数据库
	// 具体IP为服务器IP，localhost和127.0.0.1不能解析，原因未知
	videoPath := "http://10.16.43.102:8080/static/" + finalName
	figPath := "http://10.16.43.102:8080/static/" + coverPath[9:]
	if err := service.UploadVideo(uid, videoPath, figPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		})
	}
}

// PublishList all users have same publish video list
// Demand: 本机作为服务器，获取相对地址
// 6.20 此处响应正确
func PublishList(c *gin.Context) {
	userid, exist := c.Get("uid")
	if !exist {
		return
	}
	uid := userid.(int64)
	user, err := service.QueryUserById(uid)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User not exist",
			},
			VideoList: []VideoResponse{},
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: []VideoResponse{
			{
				Id:      1,
				Author:  user,
				PlayUrl: "https://www.w3schools.com/html/movie.mp4",
				//PlayUrl:       "./public/bear.mp4",
				CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
				FavoriteCount: 0,
				CommentCount:  0,
				IsFavorite:    false,
			},
		},
	})
}
