package controller

//author:fuxingyuan
import (
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type feedRes struct {
	Response
	Video_list []Video `json:"video_List,omitempty"`
	Next_time  int     `json:"next_time,omitempty"`
}

func Feed(ctx context.Context, c *app.RequestContext) {
	latest_time := c.Query("latest_time")
	videolist, err := service.Feed(ctx, c, latest_time)
	if err != nil {
		hlog.Error("视频查询错误")
		return
	}
	c.JSON(consts.StatusOK, videolist)
}
