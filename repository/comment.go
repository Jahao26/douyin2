package repository

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Comment struct {
	Cid         int64  `gorm:"column:cid;primary_key;AUTO_INCREMENT" redis:"cid"`
	Uid         int64  `gorm:"column:uid;not null" redis:"uid"`
	Vid         int64  `gorm:"column:vid;not null" redis:"vid"`
	Content     string `gorm:"column:content;type:varchar(1000);not null" redis:"content"`
	Create_date time.Time
	Delete_date gorm.DeletedAt //实现软删除
}

type CommentDAO struct {
}

var commentDao *CommentDAO
var commentOnce sync.Once

func NewCommentDao() *CommentDAO {
	commentOnce.Do(
		func() {
			commentDao = &CommentDAO{}
		})
	return commentDao
}

// AddComment 增加评论
func (*CommentDAO) AddComment(comment *Comment) (*Comment, error) {
	if err := db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// RmComment 通过cid删除评论,权限管理部分待完善
func (*CommentDAO) RmComment(commentId int64) error {
	if err := db.Where("cid = ?", commentId).Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}

// QuaryComment 通过vid查询评论
func (*CommentDAO) QuaryComment(videoId int64, commentList *[]*Comment) error {
	if commentList == nil {
		return errors.New("QueryComment commentList is nil")
	}
	err := db.Model(&Comment{}).Where("vid = ?", videoId).Find(&commentList).Error
	if err != nil {
		return err
	}
	return nil
}

// 通过cid拿到评论的uid
func (*CommentDAO) GetUidByCommentId(commentId int64) (int64, error) {
	comment := Comment{
		Cid: commentId,
	}
	err := db.Model(&Comment{}).Where("cid = ?", commentId).Find(&comment).Error
	return comment.Uid, err
}
