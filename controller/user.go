package controller

import (
	"douyin/middleware"
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var userIdSequence = int64(1)

type UserInfoPage struct {
	Id            int64
	Name          string
	FollowCount   int64
	FollowerCount int64
	IsFollow      bool
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User UserInfoPage `json:"user"`
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

//func UserInfo(c *gin.Context) {
//	// 因为登录仅获取id和name，Userinfo获取其他关注和粉丝信息
//	// userid := c.Query("user_id")
//	userid, exist := c.Get("uid")
//	if exist {
//		uid := userid.(int64)
//		user, err := service.QueryUserById(uid)
//		if err != nil {
//			c.JSON(http.StatusOK, UserResponse{Response: Response{StatusCode: 1, StatusMsg: err.Error()}})
//			return
//		}
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 0},
//			User:     user,
//		})
//	}
//}

func UserInfo(c *gin.Context) {
	// 因为登录仅获取id和name，Userinfo获取其他关注和粉丝信息
	userid := c.Query("user_id")
	uid, err := strconv.ParseInt(userid, 10, 64)
	uidget, _ := c.Get("uid") // 通过解析token得到的user_id
	fmt.Println("*****************IN USER***********")
	fmt.Println(uid)
	fmt.Println(uidget)
	fmt.Println("*****************IN USER***********")
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{Response: Response{StatusCode: 1, StatusMsg: "Id convert error"}})
	}
	_, err = service.QueryUserById(uid)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{Response: Response{StatusCode: 1, StatusMsg: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: UserInfoPage{
			Id:            10,
			Name:          "zhengjiahao",
			FollowCount:   99,
			FollowerCount: 99,
			IsFollow:      true,
		},
	})
}
