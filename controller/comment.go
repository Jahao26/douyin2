package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
// Demand: User can comment the video (need:uid, video_id, text...)
func CommentAction(c *gin.Context) {

	actionType := c.Query("action_type")

	if userid, exist := c.Get("uid"); exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			uid := userid.(int64)
			err := service.Comment(uid, text)
			if err != nil {
				panic(err)
			}
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list (before)
// Demand: Different videos hace different comment lists
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
