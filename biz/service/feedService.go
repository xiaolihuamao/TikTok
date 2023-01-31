package service

import (
	"TikTok/biz/dao"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"strconv"
	"time"
)

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author" gorm:"foreignKey:i"`
	PlayUrl       string `json:"play_url" `
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}
type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func Feed(ctx context.Context, c *app.RequestContext, latest_time string) ([]Video, error) {
	videoInfo := dao.Use(dao.Db).Video
	userInfo := dao.Use(dao.Db).User
	me, _ := strconv.ParseInt(latest_time, 10, 64)
	latesttime := time.Unix(me, 0)
	var (
		videoList []Video
		err       error
	)
	if latest_time == "" {
		err = videoInfo.WithContext(ctx).Select(userInfo.ALL, videoInfo.ALL).LeftJoin(userInfo, userInfo.UserID.EqCol(videoInfo.AuthorID)).Order(
			videoInfo.CreateTime.Desc()).Limit(30).Scan(&videoList)
		if err != nil {
			hlog.Error("查询视频数据错误")
		}
	} else {
		err = videoInfo.WithContext(ctx).LeftJoin(userInfo, userInfo.UserID.EqCol(videoInfo.AuthorID)).Where(videoInfo.CreateTime.Lt(latesttime)).Order(
			videoInfo.CreateTime.Desc()).Limit(30).Scan(&videoList)
		if err != nil {
			hlog.Error("查询视频数据错误")
		}
	}
	return videoList, err
}
