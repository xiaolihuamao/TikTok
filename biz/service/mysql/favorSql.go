package mysql

import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gen"
)

// 增加一条点赞关系
func InSertFavor(uid int64, vid int64) error {
	f := dao.Use(dao.Db).Favorite
	return f.Create(&model.Favorite{
		VideoID: vid,
		UserID:  uid,
	})
}

// 表示删除一条点赞关系
func DelFavor(uid int64, vid int64) (gen.ResultInfo, error) {
	f := dao.Use(dao.Db).Favorite
	return f.Where(f.VideoID.Eq(vid), f.UserID.Eq(uid)).Delete()
}

func GetVidByUidFromFavor(uid int64) []int64 {
	f := dao.Use(dao.Db).Favorite
	favorList, err := f.Select(f.VideoID).Where(f.UserID.Eq(uid)).Find()
	if err != nil || favorList == nil {
		hlog.Error("错误查询")
		return []int64{}
	}
	vids := make([]int64, 0, len(favorList))
	for _, favor := range favorList {
		vids = append(vids, favor.VideoID)
	}
	return vids
}
