package controller

import (
	"douyin/middleware"
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var userIdSequence = int64(1)

type UserInfoPage struct {
	// 个人主页情况
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	WorkCount     int64  `json:"work_count,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User *UserInfoPage `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//println("test!!!")
	// 查询数据库中是否有username，没有则插入数据，有则提示失败
	if id, err := service.Register(username, password); err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Register failed"}})
		return
	} else {
		token, err := middleware.GenToken(id) // 输入id，获得token
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Token process failed: " + err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if id, err := service.UserLogin(username, password); err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()}})
		return
	} else {
		token, err := middleware.GenToken(id) // 登录后从id中获取token
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    token,
		})
	}

}

func UserInfo(c *gin.Context) {
	// 因为登录仅获取id和name，Userinfo获取其他关注和粉丝信息
	// userid := c.Query("user_id")
	userid, exist := c.Get("uid")
	if exist {
		uid := userid.(int64)
		user, err := service.QueryUserById(uid)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{Response: Response{StatusCode: 1, StatusMsg: err.Error()}})
			return
		} else {
			fmt.Println(user.Id, user.Name, user.FollowerCount)
		}

		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: &UserInfoPage{
				Id:            uid,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      user.IsFollow,
				FavoriteCount: user.FavoriteCount,
				WorkCount:     user.WorkCount,
			},
		})
	}
}
