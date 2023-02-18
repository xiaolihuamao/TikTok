package service

import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	model.Comment
	User
}

func AddComment(video_id int64, comment_text string, uid int64) error {
	c := dao.Use(dao.Db).Begin()
	const shortForm = "2006-01-01"
	timeStr := time.Now().Format(shortForm)
	c.Comment.Create(&model.Comment{
		VideoID:    video_id,
		UserID:     uid,
		Content:    comment_text,
		CreateDate: timeStr,
	})
	c.Video.Where(c.Video.VideoID.Eq(video_id)).Update(c.Video.CommentCount, gorm.Expr("comment_count+?", 1))
	return c.Commit()
}

func DelComment(commentId int64, video_id int64) error {
	c := dao.Use(dao.Db).Begin()
	c.Comment.Where(c.Comment.CommentID.Eq(commentId)).Delete()
	c.Video.Where(c.Video.VideoID.Eq(video_id)).Update(c.Video.CommentCount, gorm.Expr("comment_count-?", 1))
	return c.Commit()
}

func CommentList(videoId int64) []Comment {
	c := dao.Use(dao.Db).Comment
	u := dao.Use(dao.Db).User
	var cList []Comment
	c.Select(c.VideoID, c.Content, c.UserID, c.CommentID, c.CreateDate, u.Username, u.FollowerCount, u.FollowCount).LeftJoin(u, u.UserID.EqCol(c.UserID)).Where(c.VideoID.Eq(videoId)).Scan(&cList)
	return cList
}
