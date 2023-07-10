package service

import (
	"douyin/repository"
	"time"
)

type commentFlow struct {
	uid          int64
	vid          int64
	action_type  string
	comment_text string
	comment_id   int64
}

func Comment(uid int64, vid int64, action string, text string, text_id int64) (*repository.Comment, error) {
	return newCommentFlow(uid, vid, action, text, text_id).Do()
}

func newCommentFlow(uid int64, vid int64, action string, text string, text_id int64) *commentFlow {
	return &commentFlow{uid: uid, vid: vid, action_type: action, comment_text: text, comment_id: text_id}
}

func (f *commentFlow) Do() (*repository.Comment, error) {
	comment, err := f.comment_action()
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (f *commentFlow) comment_action() (*repository.Comment, error) {
	if f.action_type == "1" {
		//当发布评论，评论保存到对应视频id用户id的评论列表里。返回当前评论内容
		newComment := &repository.Comment{
			Uid:         f.uid,
			Vid:         f.vid,
			Content:     f.comment_text,
			Create_date: time.Now(),
		}
		newComment, err := repository.NewCommentDao().AddComment(newComment)
		if err != nil {
			return nil, err
		}
		err = repository.NewVideoDao().AddVideoComment(f.vid)
		return newComment, nil

	} else {
		//当删除评论，根据评论id找到对应评论记录进行删除，返回nil
		// 进一步的，对删除者的身份进行判断
		uid, err := repository.NewCommentDao().GetUidByCommentId(f.comment_id)
		err = repository.NewVideoDao().RmVideoComment(f.vid)
		if uid == f.uid {
			err = repository.NewCommentDao().RmComment(f.comment_id)
			if err != nil {
				return nil, err
			}
		} else {
			//当id不等，不能删除评论
			return nil, nil
		}
	}
	return nil, nil
}
