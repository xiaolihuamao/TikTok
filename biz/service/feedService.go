package service

import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

type Video struct {
	model.Video
	Author User `json:"author"`
}
type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

/*
*

	返回全部视频信息的主体函数
*/
func Feed(ctx context.Context, c *app.RequestContext, latest_time time.Time, id int64) ([]Video, error) {
	var (
		videoList []Video
		err       error
	)
	videoInfo := dao.Use(dao.Db).Video
	userInfo := dao.Use(dao.Db).User
	latesttime := latest_time
	err = videoInfo.WithContext(ctx).LeftJoin(userInfo, userInfo.UserID.EqCol(videoInfo.AuthorID)).Where(videoInfo.CreateTime.Lt(latesttime)).Order(
		videoInfo.CreateTime.Desc()).Limit(3).Scan(&videoList)
	if err != nil {
		hlog.Error("查询视频数据错误")
	}
	isFavorite(id, &videoList, ctx)

	return videoList, err
}

// 判断查询出的videolist是否被登录用户点赞，也就是必须传入userid
func isFavorite(userId int64, videoList *[]Video, ctx context.Context) {
	vids := make([]int64, 0, len(*videoList))
	for _, vi := range *videoList {
		vids = append(vids, vi.VideoID)
	}
	var favorMaps []map[string]interface{}
	dao.Db.Table("favorites").Select("user_id", "video_id").Distinct().Where("video_id in ?", vids).Find(&favorMaps)
	for i, temp := range *videoList {
		for _, val := range favorMaps {
			if temp.VideoID == val["video_id"] && userId == val["user_id"] {
				(*videoList)[i].IsFavorite = 1
			}
		}
	}
}
