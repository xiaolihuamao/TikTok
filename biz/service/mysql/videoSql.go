package mysql

//author:fuxingyuan
import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"fmt"
	"gorm.io/gorm"
)

// 插入单个视频
func InsertVideo(video model.Video) error {
	v := dao.Use(dao.Db).Video
	err := v.Create(&video)

	fmt.Println(v.VideoID)
	return err
}

// 根据视频地址查询视频
func FindByPlayUrl(url string) (*model.Video, error) {
	v := dao.Use(dao.Db).Video
	return v.Select(v.ALL).Where(v.PlayURL.Eq(url)).First()
}

// 根据作者id查询视频信息
func FindByAuthor(authorId int64) ([]*model.Video, error) {
	v := dao.Use(dao.Db).Video
	return v.Select(v.ALL).Where(v.AuthorID.Eq(authorId)).Find()
}

// 更新点赞数
func SetFavorCountById(vid int64, favorCount int64) (int64, error) {
	v := dao.Use(dao.Db).Video
	info, err := v.Select(v.FavoriteCount).Where(v.VideoID.Eq(vid)).Update(v.FavoriteCount, gorm.Expr("favorite_count+?", favorCount))
	return info.RowsAffected, err
}

// 根据视频id查询作者
func FindAthorByVid(vid int64) (int64, error) {
	v := dao.Use(dao.Db).Video
	info, err := v.Select(v.AuthorID).Where(v.VideoID.Eq(vid)).First()
	return info.AuthorID, err
}
