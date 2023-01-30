package controller

//author:fuxingyuan
import (
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type feedRes struct {
	Response
	Video_list []Video `json:"video_List,omitempty"`
	Next_time  int     `json:"next_time,omitempty"`
}

func Feed(ctx context.Context, c *app.RequestContext) {
	latest_time := c.Query("latest_time")
	service.Feed(ctx, c, latest_time)
	c.JSON(consts.StatusOK, utils.H{
		"message": "pong",
	})
}
