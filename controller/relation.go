package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction 点击关注按钮，添加到follow数据库和被关注人的follower数据库
func RelationAction(c *gin.Context) {
	token := c.Query("token")           // 当前用户的token
	to_user_id := c.Query("to_user_id") // 被关注人id
	intnum, _ := strconv.Atoi(to_user_id)
	touserid := int64(intnum)

	if user, exist := usersLoginInfo[token]; exist {
		newRela := Relation{
			Follower_id:  touserid,
			Following_id: user.Id,
			Create_at:    time.Now(),
		}
		db.Create(&newRela) // 插入新的数据
		// 修改

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
