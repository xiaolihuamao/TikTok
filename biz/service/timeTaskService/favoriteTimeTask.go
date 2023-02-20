package timeTaskService

import (
	redisUtil "TikTok/biz/mw/redis"
	"TikTok/biz/service/mysql"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-co-op/gocron"
	"strconv"
	"time"
)

func SyncFavorDb() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	s := gocron.NewScheduler(timezone)
	// 每秒执行一次
	s.Every(30).Seconds().Do(func() {
		go func() {
			isFavorVids, err := redisUtil.Rdb.SMembers("isLiked_video_set").Result()
			fmt.Println("任务在执行")
			if err != nil || isFavorVids == nil {
				hlog.Error("定时任务失败")
				return
			}
			for _, vid := range isFavorVids {
				var (
					err        error
					favorCount int64
					vidInt     int64
				)
				vidInt, err = strconv.ParseInt(vid, 0, 64)
				if vidInt <= 0 {
					continue
				}
				key := "videoLike_count_" + vid

				favorCount, err = redisUtil.Rdb.Get(key).Int64()

				if err != nil {
					return
				}
				//根据vid取出给uid点赞过的对应的uid
				uids, err := redisUtil.Rdb.SMembers("videoLiked_users_" + vid).Result()
				if err != nil || uids == nil {
					return
				}
				for _, uid := range uids {
					uidInt, _ := strconv.ParseInt(uid, 0, 64)
					if uidInt <= 0 {
						continue
					}
					relationKey := "like" + vid + "::" + uid
					code, err := redisUtil.Rdb.Get(relationKey).Int()
					if err != nil {
						return
					}
					if code == 0 {
						//表示是取消
						_, err := mysql.DelFavor(uidInt, vidInt)
						if err != nil {
							return
						}
					} else if code == 1 {
						//表示是点赞
						err := mysql.InSertFavor(uidInt, vidInt)
						if err != nil {
							return
						}
					}
					//刷盘后删除关系键
					redisUtil.Rdb.Del(relationKey)
					//移除元素
					redisUtil.Rdb.SRem("videoLiked_users_"+vid, uid)
					//移除元素
					redisUtil.Rdb.SRem("userlike_videos_"+uid, vid)
				}
				_, err = mysql.SetFavorCountById(vidInt, favorCount)
				if err != nil {
					hlog.Error("数据库同步错误")
					return
				} else {
					//刷盘后删除对于vid点赞统计的key
					redisUtil.Rdb.Del(key)
					//删除被点赞video集合中对应vid
					redisUtil.Rdb.SRem("isLiked_video_set", vid)
				}
			}
		}()
	})
	s.StartAsync()
}
