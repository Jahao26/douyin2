package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []*service.NewComment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment *service.NewComment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
// Demand: User can comment the video (need:uid, video_id, text...)
func CommentAction(c *gin.Context) {

	actionType := c.Query("action_type")
	videoId := c.Query("video_id")
	commentId := c.Query("comment_text")

	if userid, exist := c.Get("uid"); exist {
		text := c.Query("comment_text")
		uid := userid.(int64)
		vid, _ := strconv.ParseInt(videoId, 10, 64)
		cid, _ := strconv.ParseInt(commentId, 10, 64)
		comment, err := service.Comment(uid, vid, actionType, text, cid)
		if err != nil {
			panic(err)
		}
		// 通过uid获得userinfo
		user, err := service.UserInfo(comment.Uid)

		c.JSON(http.StatusOK, CommentActionResponse{
			Comment: &service.NewComment{
				Id:         comment.Cid,
				User:       user,
				Content:    comment.Content,
				CreateDate: comment.Create_date.Format("2006-01-02 15:04:05"),
			},
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list (before)
// Demand: Different videos hace different comment lists
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")
	if _, exist := c.Get("uid"); exist {
		//uid := userid.(int64)
		vid, _ := strconv.ParseInt(videoId, 10, 64)
		commentlist, _ := service.QuaryCommentList(vid)

		c.JSON(http.StatusOK, CommentListResponse{
			Response:    Response{StatusCode: 0},
			CommentList: commentlist.Comments,
		})

	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}
