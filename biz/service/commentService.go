package service

//author:zhangshuo
import (
	"TikTok/biz/controller"
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"log"
)

/*
	评论模块自己request实现的方法：
*/

// Send
// 1.发表评论，传进来评论的基本信息，返回保存是否成功的状态描述
func Send(comment model.Comment) (controller.CommentInfo, error) {
	log.Println("CommentService-Send: running") //函数已运行
	//数据准备
	var commentInfo model.Comment
	commentInfo.VideoID = comment.VideoID       //评论视频id传入
	commentInfo.UserID = comment.UserID         //评论用户id传入
	commentInfo.Content = comment.Content       //评论内容传入
	commentInfo.CreateDate = comment.CreateDate //评论时间
	//commentInfo.Cancel = conf.ValidComment      //评论状态，0，有效

	//1.评论信息存储：
	commentRtn, err := dao.InsertComment(commentInfo)
	if err != nil {
		return controller.CommentInfo{}, err
	}

	//2.查询用户信息
	userData, err2 := dao.GetUserByIdWithCurId(comment.UserID, comment.UserID)
	if err2 != nil {
		return controller.CommentInfo{}, err2
	}

	//3.拼接,
	commentData := controller.CommentInfo{
		Id:         commentRtn.CommentID,
		User:       userData,
		Content:    commentRtn.Content,
		CreateDate: commentRtn.CreateDate,
	}

	//返回结果
	return commentData, nil
}

// DelComment
// 2.删除评论，传入评论id即可，返回错误状态信息,开发中！！！
func DelComment(commentId int64) error {
	return nil
}

// GetList
// 3.查看评论列表-返回评论list-在controller层再封装外层的状态信息，开发中！！！
func GetList(videoId int64, userId int64) ([]controller.CommentInfo, error) {
	return []controller.CommentInfo{}, nil
}
