package service

import "douyin/repository"

type NewComment struct {
	Id         int64         `json:"id,omitempty"`
	User       *UserInfoPage `json:"user"`
	Content    string        `json:"content,omitempty"`
	CreateDate string        `json:"create_date,omitempty"`
}

type CommentList struct {
	Comments []*NewComment
}

type QuaryCommentListFlow struct {
	vid         int64
	comments    []*repository.Comment
	commentlist *CommentList
}

func QuaryCommentList(vid int64) (*CommentList, error) {
	return newQuaryCommentListFlow(vid).Do()
}

func newQuaryCommentListFlow(vid int64) *QuaryCommentListFlow {
	return &QuaryCommentListFlow{vid: vid}
}

func (q *QuaryCommentListFlow) Do() (*CommentList, error) {
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.commentlist, nil
}

func (q *QuaryCommentListFlow) packData() error {
	//先获得原始评论列表
	err := repository.NewCommentDao().QuaryComment(q.vid, &q.comments)
	if err != nil {
		return err
	}
	//对评论列表进行转换
	newcommentlist := make([]*NewComment, len(q.comments))
	for i := range q.comments {
		user, _ := UserInfo(q.comments[i].Uid)
		newcomment := NewComment{
			Id:         q.comments[i].Cid,
			User:       user,
			Content:    q.comments[i].Content,
			CreateDate: q.comments[i].Create_date.Format("2006-01-02 15:04:05"),
		}
		newcommentlist[i] = &newcomment
	}
	q.commentlist = &CommentList{Comments: newcommentlist}
	return nil
}
