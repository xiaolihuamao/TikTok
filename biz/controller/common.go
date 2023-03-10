package controller

import "time"

//author:fuxingyuan
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64     `json:"id,omitempty"`
	Author        User      `json:"author"`
	PlayUrl       string    `json:"play_url" `
	CoverUrl      string    `json:"cover_url"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	IsFavorite    bool      `json:"is_favorite"`
	Title         string    `json:"title,omitempty"`
	Create_time   time.Time `json:"create_time"`
}

type CommentInfo struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type User struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	FollowCount      int64  `json:"follow_count"`
	FollowerCount    int64  `json:"follower_count"`
	IsFollow         bool   `json:"is_follow"`
	Favorite_count   int64  `json:"favorite_count"`
	Avatar           string `json:"avatar"`           //用户头像
	Background_image string `json:"background_image"` //用户顶部页大图
	Signature        string `json:"signature"`        //个人简介
	Total_favorited  string `json:"total_favorited"`  //获赞数量
	Work_count       int64  `json:"work_count"`       //作品数
}
