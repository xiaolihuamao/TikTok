package service

import (
	"TikTok/biz/controller"
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"time"
)

func Feed(ctx context.Context, c *app.RequestContext, latest_time string) ([]controller.Video, error) {
	videoInfo := dao.Use(dao.Db).Video
	userInfo := dao.Use(dao.Db).User
	me, _ := strconv.ParseInt(latest_time, 10, 64)
	latesttime := time.Unix(me, 0)
	var (
		videoList []*model.Video
		err       error
	)
	if latest_time == "" {
		videoList, err = videoInfo.WithContext(ctx).Joins(on()).Order(
			videoInfo.CreateTime.Desc()).Limit(30).Find()
		if err != nil {
			hlog.Error("查询视频数据错误")
		}
	} else {
		videoList, err = videoInfo.WithContext(ctx).Where(videoInfo.CreateTime.Lt(latesttime)).Order(videoInfo.CreateTime.Desc()).Limit(30).Find()
		if err != nil {
			hlog.Error("查询视频数据错误")
		}
	}
}
