package dao

//author:zhangshuo

import (
	"TikTok/conf"
	"errors"
	"log"
	"time"
)

// CommentData 对应数据库Comment表结构的结构体
type CommentData struct {
	Id          int64     //评论id
	UserId      int64     //评论用户id
	VideoId     int64     //视频id
	CommentText string    //评论内容
	CreateDate  time.Time //评论发布的日期mm-dd
	Cancel      int32     //取消评论为1，发布评论为0
}

// TableUser 对应数据库User表结构的结构体
type TableUser struct {
	Id       int64
	Name     string
	Password string
}

// InsertComment
// 将评论插入数据库
func InsertComment(comment CommentData) (CommentData, error) {
	log.Println("CommentDao-InsertComment: running") //函数已运行
	//数据库中插入一条评论信息
	err := Db.Model(CommentData{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") //函数返回提示错误信息
		return CommentData{}, errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success") //函数执行成功，返回正确信息
	return comment, nil
}

// DeleteComment
// 将评论从数据库中删除
func DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running") //函数已运行
	var commentInfo CommentData
	//先查询是否有此评论
	result := Db.Model(CommentData{}).Where(map[string]interface{}{"id": id, "cancel": conf.ValidComment}).First(&commentInfo)
	if result.RowsAffected == 0 { //查询到此评论数量为0则返回无此评论
		log.Println("CommentDao-DeleteComment: return del comment is not exist") //函数返回提示错误信息
		return errors.New("del comment is not exist")
	}
	//数据库中删除评论-更新评论状态为-1
	err := Db.Model(CommentData{}).Where("id = ?", id).Update("cancel", conf.InvalidComment).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") //函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") //函数执行成功，返回正确信息
	return nil
}

// GetCommentList
// 根据视频id查询所属评论全部列表信息
func GetCommentList(videoId int64) ([]CommentData, error) {
	log.Println("CommentDao-GetCommentList: running") //函数已运行
	//数据库中查询评论信息list
	var commentList []CommentData
	result := Db.Model(CommentData{}).Where(map[string]interface{}{"video_id": videoId, "cancel": conf.ValidComment}).
		Order("create_date desc").Find(&commentList)
	//若此视频没有评论信息，返回空列表，不报错
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments") //函数返回提示无评论
		return nil, nil
	}
	//若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed") //函数返回提示获取评论错误
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success") //函数执行成功，返回正确信息
	return commentList, nil
}

// GetTableUserById
// 根据user_id获得TableUser对象
func GetTableUserById(id int64) (TableUser, error) {
	tableUser := TableUser{}
	if err := Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}
