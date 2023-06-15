package service

type commentFlow struct {
}

func Comment(uid int64, text string) error {
	return newCommentFlow(uid, text).Do()
}

func newCommentFlow(uid int64, text string) *commentFlow {
	return &commentFlow{}
}

func (f *commentFlow) Do() error {
	return nil
}

func (f *commentFlow) comment_action() {

}
