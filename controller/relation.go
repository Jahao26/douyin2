package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction 点击关注按钮，添加到follow数据库和被关注人的follower数据库
func RelationAction(c *gin.Context) {
	touserid := c.Query("to_user_id") // 被关注人id
	intnum, _ := strconv.Atoi(touserid)
	to_user_id := int64(intnum)

	if userid, exist := c.Get("uid"); exist {
		uid := userid.(int64)
		if err := service.RelationAction(uid, to_user_id); err != nil {
			panic(err)
		}
		//uid := userid.(int64)
		//newRela := Relation{
		//	Follower_id:  to_user_id,
		//	Following_id: uid,
		//	Create_at:    time.Now(),
		//}

		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
// demand: 不同用户需要有不用的关注列表
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FollowerList all users have same follower list
// demand：不同用户需要有不同的粉丝列表
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FriendList all users have same friend list
// demand：如果用户互关，他们就是朋友，展示朋友列表
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
