package service

import (
	"TikTok/biz/dao"
	redisUtil "TikTok/biz/mw/redis"
	"TikTok/biz/service/mysql"
	"context"
	"strconv"
)

// 点赞Service
func AddLike(uid int64, vid int64, ctx context.Context) error {
	vidStr := strconv.Itoa(int(vid))
	uidStr := strconv.Itoa(int(uid))
	pipe := redisUtil.Rdb.Pipeline()
	//将所有被点赞的videoid存入一个set集合
	pipe.SAdd("isLiked_video_set", vid)
	//构建一个以videoId为key,存放所有对该video点过赞的用户id
	pipe.SAdd("videoLiked_users_"+vidStr, uid)
	//构建一个以uid为key,存放所有点赞的视频videoid
	pipe.SAdd("userlike_videos_"+uidStr, vid)
	//incr统计点赞数量，之所以不scard上面的video-(users)集合是因为取消点赞时可能已经刷盘，所以redis要体现出负数点赞
	pipe.Incr("videoLike_count_" + vidStr)
	//set 集合，用vid和uid复合作为key,value 1 表示已经点赞，取消可以设为0
	pipe.Set("like"+vidStr+"::"+uidStr, 1, 0)
	_, err := pipe.Exec()
	return err
}

// 取消点赞的Service
func CancelLike(uid int64, vid int64, ctx context.Context) error {
	vidStr := strconv.Itoa(int(vid))
	uidStr := strconv.Itoa(int(uid))
	pipe := redisUtil.Rdb.Pipeline()
	//将所有被取消点赞的videoid存入一个set集合，与点赞同理，这个集合更准确是发生过赞关系的都要收集
	pipe.SAdd("isLiked_video_set", vid)
	//构建一个以videoId为key,存放所有对该video取消过赞的用户id
	pipe.SAdd("videoLiked_users_"+vidStr, uid)
	//构建一个以uid为key,存放所有取消赞的视频videoid
	pipe.SAdd("userlike_videos_"+uidStr, vid)
	pipe.Decr("videoLike_count_" + vidStr)
	pipe.Set("like"+vidStr+"::"+uidStr, 0, 0)
	_, err := pipe.Exec()
	return err
}

// 喜欢的视频列表
func FavorList(uid int64, ctx context.Context) []Video {
	var vids []int64
	uidStr := strconv.Itoa(int(uid))
	vids = mysql.GetVidByUidFromFavor(uid)
	vidStrs, err := redisUtil.Rdb.SMembers("userlike_videos_" + uidStr).Result()
	if err == nil && vidStrs != nil {
		for _, vidstr := range vidStrs {
			key := "like" + vidstr + "::" + uidStr
			code, err := redisUtil.Rdb.Get(key).Int64()
			if err == nil && code == 1 {
				vid, _ := strconv.ParseInt(vidstr, 0, 64)
				vids = append(vids, vid)
			}
		}
	}
	var videoList []Video
	dao.Db.Table("videos").Select("videos.*", "users.*").Joins("left join users on users.user_id=videos.author_id").Find(&videoList, vids)
	createVideo(uid, &videoList, ctx)
	return videoList
}
