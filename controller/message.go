package controller

import (
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//var tempChat = map[string][]Message{}
//
//var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []*service.NewMassage `json:"message_list,omitempty"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	toUid := c.Query("to_user_id")
	actionType := c.Query("action_type")
	content := c.Query("content")

	if userId, exist := c.Get("uid"); exist {
		uid := userId.(int64)
		to_uid, _ := strconv.ParseInt(toUid, 10, 64)
		// 发送消息时，保存消息记录，返回响应
		if err := service.MassageAction(uid, to_uid, actionType, content); err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	toUid := c.Query("to_user_id")

	if userid, exist := c.Get("uid"); exist {
		uid := userid.(int64)
		to_uid, _ := strconv.ParseInt(toUid, 10, 64)

		massagelist, _ := service.QuaryMassageList(uid, to_uid)
		//fmt.Println("****")
		//fmt.Println(uid, to_uid)
		//for i := range massagelist.Massages {
		//	fmt.Println(massagelist.Massages[i])
		//}

		c.JSON(http.StatusOK, ChatResponse{
			Response:    Response{StatusCode: 0},
			MessageList: massagelist.Massages,
		})

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
