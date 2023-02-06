package service

//author :fuxingyuan
import (
	"TikTok/biz/dao"
	"TikTok/biz/model"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"sync"
	"time"
)

/*
service封装的videolist，将model层的Vedio,User取出来。
*/
type Video struct {
	model.Video
	model.User
	Is_follow bool `json:"is_follow"`
}

/*
		返回全部视频信息的主体函数
	     latest_time 最新时间
	     id 当前登录用户id。
*/
func Feed(ctx context.Context, c *app.RequestContext, latest_time time.Time, id int64) ([]Video, error) {
	var (
		videoList []Video
		err       error
	)
	v := dao.Use(dao.Db).Video
	u := dao.Use(dao.Db).User
	latesttime := latest_time
	//查询最新视频流及作者
	err = v.WithContext(ctx).Select(v.VideoID,
		v.AuthorID, v.CreateTime, v.Title, v.CoverURL, v.PlayURL, v.CommentCount, v.FavoriteCount,
		u.UserID, u.Username, u.FollowerCount, u.FollowCount,
	).LeftJoin(u, u.UserID.EqCol(v.AuthorID)).Where(v.CreateTime.Lt(latesttime)).Order(
		v.CreateTime.Desc()).Limit(1).Scan(&videoList)
	if err != nil {
		hlog.Error("查询视频数据错误")
	}
	//登录用户执行此步操作，判断是否isFavorite
	if id >= 0 {
		creatVideoList(id, &videoList, ctx)
	}
	return videoList, err
}

func creatVideoList(id int64, videoList *[]Video, ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)
	//定义判断是否点过赞的函数
	var isFavorite func(userId int64, videoList *[]Video, ctx context.Context)
	isFavorite = func(userId int64, videoList *[]Video, ctx context.Context) {
		vids := make([]int64, 0, len(*videoList))
		for _, vi := range *videoList {
			vids = append(vids, vi.VideoID)
		} //将所有列表中的videoID取出作为集合
		var favorMaps []map[string]interface{}
		//在favorite表中查询所有user_id video_id map
		dao.Db.Table("favorites").Select("user_id", "video_id").Distinct().Where("video_id in ?", vids).Find(&favorMaps)
		//遍历videoList,遍历favorMap，将videoList中VideoID=map["video_id"],且此map下map["user_id"]=userId的video.isFavorite赋值1，表示存在点赞关系。
		for i, temp := range *videoList {
			for _, val := range favorMaps {
				if temp.VideoID == val["video_id"] && userId == val["user_id"] {
					(*videoList)[i].IsFavorite = 1
				}
			}
		}
		wg.Done()
	}
	//定义测试函数，后期可改造为从缓存取出点赞数
	var addFavorNumFromCache func(userId int64, videoList *[]Video, ctx context.Context)
	addFavorNumFromCache = func(userId int64, videoList *[]Video, ctx context.Context) {
		for i, temp := range *videoList {
			if temp.VideoID%2 == 0 {
				(*videoList)[i].FavoriteCount = 100
			} else {
				(*videoList)[i].FavoriteCount = 1000
			}
		}
		wg.Done()
	}
	go isFavorite(id, videoList, ctx)
	go addFavorNumFromCache(id, videoList, ctx)
	wg.Wait()
}

/*
	判断查询出的videolist是否被登录用户点赞，也就是必须传入userid

userId登录的用户id
videoList 待处理的Video列表
*/
/*func isFavorite(userId int64, videoList *[]Video, ctx context.Context) {
	vids := make([]int64, 0, len(*videoList))
	for _, vi := range *videoList {
		vids = append(vids, vi.VideoID)
	} //将所有列表中的videoID取出作为集合
	var favorMaps []map[string]interface{}
	//在favorite表中查询所有user_id video_id map
	dao.Db.Table("favorites").Select("user_id", "video_id").Distinct().Where("video_id in ?", vids).Find(&favorMaps)
	//遍历videoList,遍历favorMap，将videoList中VideoID=map["video_id"],且此map下map["user_id"]=userId的video.isFavorite赋值1，表示存在点赞关系。
	for i, temp := range *videoList {
		for _, val := range favorMaps {
			if temp.VideoID == val["video_id"] && userId == val["user_id"] {
				(*videoList)[i].IsFavorite = 1
			}
		}
	}
}*/
