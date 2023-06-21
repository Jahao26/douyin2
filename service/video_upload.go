package service

import (
	"bytes"
	"douyin/repository"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

type uploadFlow struct {
	id        int64
	videopath string
	coverpath string
}

func UploadVideo(uid int64, vpath string, cpath string) error {
	return NewUploadFlow(uid, vpath, cpath).Do()
}

func GetCoverimage(videoPath string, imageName string) (string, error) {
	// 输出封面图片的路径
	snapshotPath := "./public/" + imageName
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(videoPath).Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("缩略图解码失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		log.Fatal("缩略图保存失败：", err)
		return "", err
	}

	imgPath := snapshotPath + ".png"

	return imgPath, nil
}

func NewUploadFlow(uid int64, vpath string, cpath string) *uploadFlow {
	return &uploadFlow{
		id:        uid,
		videopath: vpath,
		coverpath: cpath,
	}
}

func (f *uploadFlow) Do() error {
	if err := f.uploadVideo(); err != nil {
		return errors.New("Error in upload")
	}
	return nil
}

func (f *uploadFlow) uploadVideo() error {
	// 查询uid对应的姓名
	user, err := repository.NewUserDao().QueryById(f.id)
	if err != nil {
		return err
	}

	// 将视频信息放入数据库
	newVideo := repository.Video{
		Id:            0,
		Uid:           user.Id,
		PlayUrl:       f.videopath,
		CoverUrl:      f.coverpath,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	}

	if err := repository.NewVideoDao().AddVideo(&newVideo); err != nil {
		return err
	}
	return nil
}
