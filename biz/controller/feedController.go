package controller

//author:fuxingyuan
import (
	mw "TikTok/biz/mw/jwt"
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"time"
)

type feedRes struct {
	Response
	Video_list []Video     `json:"video_List,omitempty"`
	Next_time  interface{} `json:"next_time,omitempty"`
}

func Feed(ctx context.Context, c *app.RequestContext) {
	latest_time := c.Query("latest_time")
	var lastTime time.Time
	if latest_time != "" && latest_time != "0" {
		me, _ := strconv.ParseInt(latest_time, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	unknown, exists := c.Get("id")
	var uid int64
	if !exists {
		uid = -1
		hlog.Info("未登录")
	}
	userDemo := unknown.(*mw.UserDemo)
	uid = userDemo.Uid
	videolist, err := service.Feed(ctx, c, lastTime, uid)
	var videoRes = []Video{}
	copyVideo(&videolist, &videoRes)
	if err != nil {
		hlog.Error("视频查询错误")
		c.JSON(consts.StatusInternalServerError, feedRes{
			Response:   Response{StatusCode: 1, StatusMsg: "查询视频错误"},
			Video_list: nil,
			Next_time:  nil,
		})
		return
	}
	var next_time int64
	if len(videoRes) != 0 {
		next_time = videoRes[len(videoRes)-1].Create_time.Unix()
	} else {
		next_time = time.Now().Unix()
	}
	c.JSON(consts.StatusOK, feedRes{
		Response:   Response{StatusCode: 0, StatusMsg: "success"},
		Video_list: videoRes,
		Next_time:  next_time,
	})
}

// 将service数据复刻到controller封装类
func copyVideo(v1 *[]service.Video, v2 *[]Video) {
	for _, temp := range *v1 {
		user := User{
			Id:            temp.Author.Id,
			Name:          temp.Author.Name,
			FollowCount:   temp.Author.FollowCount,
			FollowerCount: temp.Author.FollowerCount,
			IsFollow:      temp.Author.IsFollow,
		}
		var isFavorite bool
		if temp.IsFavorite != 0 {
			isFavorite = true
		}
		video := Video{
			Id:            temp.VideoID,
			Author:        user,
			PlayUrl:       temp.PlayURL,
			CoverUrl:      temp.CoverURL,
			FavoriteCount: temp.FavoriteCount,
			CommentCount:  temp.CommentCount,
			IsFavorite:    isFavorite,
			Title:         temp.Title,
			Create_time:   temp.CreateTime,
		}
		*v2 = append(*v2, video)
	}
}
