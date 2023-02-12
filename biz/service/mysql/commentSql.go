package mysql

import (
	"TikTok/biz/dao"
	"gorm.io/gen"
)

func DelCommentByCid(cid int64) (info gen.ResultInfo, err error) {
	c := dao.Use(dao.Db).Comment
	return c.Where(c.CommentID.Eq(cid)).Delete()
}
