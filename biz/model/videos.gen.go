// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVideo = "videos"

// Video mapped from table <videos>
type Video struct {
	VideoID       int64     `gorm:"column:video_id;type:bigint;primaryKey;autoIncrement:true" json:"video_id"`
	AuthorID      int64     `gorm:"column:author_id;type:bigint;not null" json:"author_id"`
	Title         string    `gorm:"column:title;type:varchar(300)" json:"title"`
	PlayURL       string    `gorm:"column:play_url;type:varchar(500)" json:"play_url"`
	CoverURL      string    `gorm:"column:cover_url;type:varchar(500)" json:"cover_url"`
	CreateTime    time.Time `gorm:"column:create_time;type:int unsigned;autoCreateTime" json:"create_time"`
	FavoriteCount int64     `gorm:"column:favorite_count;type:bigint;not null" json:"favorite_count"`
	CommentCount  int64     `gorm:"column:comment_count;type:bigint;not null" json:"comment_count"`
	IsFavorite    int64     `gorm:"column:is_favorite;type:tinyint;not null" json:"is_favorite"`
}

// TableName Video's table name
func (*Video) TableName() string {
	return TableNameVideo
}
