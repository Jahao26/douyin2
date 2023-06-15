package controller

import (
	"douyin/middleware"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

//type UserResponse struct { // ???
//	Response
//	User *service.User `json:"user"`
//}

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

	//token := username + password
	//var userlog User
	//db.Where("name=?", username).First(&userlog)
	//if userlog.Id != 0 { // 登录时用户存在在数据库内
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userlog.Id,
	//		Token:    token,
	//	})
	//	//usersLoginInfo[token] = userlog
	//} else {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	//	})
	//}
}

//// 验证用户是否存在
//func UserInfo(c *gin.Context) {
//	token := c.Query("token")
//	if user, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 0},
//			User:     user,
//		})
//
//	} else {
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
//		})
//	}
//}
