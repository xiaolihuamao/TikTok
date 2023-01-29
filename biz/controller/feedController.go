package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Feed(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{
		"message": "pong",
	})
}
