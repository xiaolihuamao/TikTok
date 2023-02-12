package service

import (
	"TikTok/biz/model"
	"TikTok/biz/service/mysql"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	uuid "github.com/satori/go.uuid"
	"mime/multipart"
	"os/exec"
	"runtime"
	"time"
)

func Publish(uid int64, file *multipart.FileHeader, c *app.RequestContext, title string, ctx context.Context) error {
	var err error
	uuid := uuid.NewV1()
	//拼接url前缀，测试时需要修改这个
	urlPrefix := "http://192.168.137.1:8081/upload/"
	//拼接url,这个是存进数据库的
	url := fmt.Sprintf("%s%v.mp4", urlPrefix, uuid)
	videoExist, err := mysql.FindByPlayUrl(url)
	if videoExist != nil {
		hlog.Error("视频重复投稿！")
		return err
	}
	//拼接保存地址
	saveRoot := fmt.Sprintf("./file/upload/%v.mp4", uuid)
	//执行保存
	err = c.SaveUploadedFile(file, saveRoot)

	//获取当前文件绝对路径
	_, path, _, _ := runtime.Caller(0)
	//推算视频图片文件静态绝对路径
	path = path[0:len(path)-len("/biz/service/publishservice.go")] + "/file/upload/"
	//拼接视频路径
	vpath := path + fmt.Sprintf("%v.mp4", uuid)
	//拼接封面保存路径
	ppath := fmt.Sprintf("%s%v.jpg", path, uuid)
	//执行ffmpeg命令，vpath表示取到视频的绝对路径，ppath表示保存图片的路径前缀
	cmd := exec.CommandContext(ctx, "cmd", "/C", "ffmpeg -ss 00:00:01 -i "+vpath+" -frames:v 1 "+ppath)
	err = cmd.Run()
	if err != nil {
		hlog.Error("截取封面错误")
		return err
	}
	//拼接图片url路径
	urlPic := urlPrefix + fmt.Sprintf("%v.jpg", uuid)

	video := model.Video{
		AuthorID:      uid,
		Title:         title,
		PlayURL:       url,
		CoverURL:      urlPic,
		CreateTime:    time.Now(),
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    0,
	}
	fmt.Println(video)
	err = mysql.InsertVideo(video)
	return err
}

// 查询用户发布的视频
func PublishList(userId int64, ctx context.Context) []Video {
	userInfo, err := mysql.GetUserById(userId)

	if err != nil {
		userInfo = new(model.User)
	}
	videoInfo, err := mysql.FindByAuthor(userId)
	var videoList = make([]Video, 0, len(videoInfo))
	for _, temp := range videoInfo {
		video := Video{
			Video: *temp,
			User: User{
				User:      *userInfo,
				Is_follow: true,
			},
		}
		videoList = append(videoList, video)
	}
	createVideo(userId, &videoList, ctx)
	return videoList
}
