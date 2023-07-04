package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []*service.UserInfoPage `json:"user_list"`
}

// RelationAction 点击关注按钮，添加到follow数据库和被关注人的follower数据库
func RelationAction(c *gin.Context) {
	touserid := c.Query("to_user_id") // 被关注人id
	to_user_id, _ := strconv.ParseInt(touserid, 10, 64)
	action_type := c.Query("action_type")

	if userid, exist := c.Get("uid"); exist {
		uid := userid.(int64)
		if err := service.RelationAction(uid, to_user_id, action_type); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
// demand: 不同用户需要有不用的关注列表
func FollowList(c *gin.Context) {
	userid, exist := c.Get("uid")
	if !exist {
		return
	}
	uid := userid.(int64)
	_, err := service.UserInfo(uid)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User not exist",
			},
			VideoList: []*service.VideoResponse{},
		})
	}
	// 通过uid获得当前用户的关注列表
	userlist, err := service.FollowList(uid)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userlist.Users,
	})
}

// FollowerList all users have same follower list
// demand：不同用户需要有不同的粉丝列表
func FollowerList(c *gin.Context) {
	userid, exist := c.Get("uid")
	if !exist {
		return
	}
	uid := userid.(int64)
	_, err := service.UserInfo(uid)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User not exist",
			},
			VideoList: []*service.VideoResponse{},
		})
	}
	// 通过uid获得当前用户的关注列表
	userlist, err := service.FollowerList(uid)

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userlist.Users,
	})
}

// FriendList all users have same friend list
// demand：如果用户互关，他们就是朋友，展示朋友列表
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: nil,
	})
}
