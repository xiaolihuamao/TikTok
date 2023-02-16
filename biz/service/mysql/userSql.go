package mysql

import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
)

// 根据用户id查询用户信息
func GetUserById(uid int64) (*model.User, error) {
	u := dao.Use(dao.Db).User
	return u.Select(u.ALL).Where(u.UserID.Eq(uid)).First()
}
