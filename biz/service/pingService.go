package service

//@author fuxingyuan
import (
	"TikTok/biz/dao"
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Ping(ctx context.Context) {
	user := dao.Use(dao.Db).User
	_, err := user.WithContext(ctx).First()
	if err != nil {
		hlog.Error("查询不到数据")
	}
}
