package mysql

//author fuxingyuan
import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
)

// 根据用户id查询用户信息
func GetUserById(uid int64) (*model.User, error) {
	u := dao.Use(dao.Db).User
	return u.Select(u.ALL).Where(u.UserID.Eq(uid)).First()
}

// 注册插入一条用户数据并返回主键
func InsertUser(user *model.User) (int64, error) {
	u := dao.Use(dao.Db).User
	err := u.Create(user)
	return user.UserID, err
}
