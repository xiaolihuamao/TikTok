package dao

//author:zhangshuo

import (
	"TikTok/biz/controller"
	"TikTok/biz/model"
	"errors"
	"log"
)

// InsertComment
// 将评论插入数据库
func InsertComment(comment model.Comment) (model.Comment, error) {
	log.Println("CommentDao-InsertComment: running") //函数已运行
	//数据库中插入一条评论信息
	err := Db.Model(model.Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") //函数返回提示错误信息
		return model.Comment{}, errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success") //函数执行成功，返回正确信息
	return comment, nil
}

// GetUserByIdWithCurId
// 由登录用户id获取用户信息，开发中！！！
func GetUserByIdWithCurId(UserID1 int64, UserID2 int64) (controller.User, error) {
	return controller.User{}, nil
}
