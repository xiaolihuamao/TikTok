// Code generated by hertz generator.
// @author fuxingyuan
package controller

import (
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Ping .
func Ping(ctx context.Context, c *app.RequestContext) {
	service.Ping(ctx)
	c.JSON(consts.StatusOK, utils.H{
		"message":    "pong",
		"statuscode": 0,
	})
}
