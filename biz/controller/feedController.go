package controller

//author:fuxingyuan
import (
	mw "TikTok/biz/mw/jwt"
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
	"time"
)

/*
*
feed 流封装返回类
Response:controller/common中封装的统一返回类
Video:controller/common中封装的统一视频列表返回类
Next_time:下一次的时间戳
*/
type feedRes struct {
	Response
	Video_list []Video     `json:"video_list,omitempty"`
	Next_time  interface{} `json:"next_time,omitempty"`
}

/*
*

	/douyin/feed/?token="adfajgfajgfuafgajfg"&latest_time="199999999" Get
*/
func Feed(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token") //获取url中的token
	var uid interface{}       //接收token解析出的uid
	//token不为空，解析token
	if token != "" {
		claims, err := mw.AuthMiddleware.GetClaimsFromJWT(ctx, c) //解析token,取出claims map
		if err != nil {
			hlog.Error("token解析错误，请使用正确的token")
			uid = float64(-1)
		}
		//取出登录后返回的token中保存的uid---(interface{}/float64)
		uid = claims["id"]
		//若取不出uid,说明token错误或过期，给uid赋值float64(-1)
		if uid == nil {
			uid = float64(-1)
			hlog.Error("token解析错误，请使用正确的token")
		}
		//token为空，给uid赋值float(-1)
	} else {
		uid = float64(-1)
		hlog.Info("未登录用户")
	}
	uidf := uid.(float64)
	uidInt := int64(uidf)                 //uid interface{}/float64--->int64 方便传参
	latest_time := c.Query("latest_time") //取出时间戳
	var lastTime time.Time                //定义Time类型的lastTime
	//对时间戳判空，若为空，取当前时间，不为空，将取出的时间戳转time类型
	if latest_time != "" && latest_time != "0" {
		me, _ := strconv.ParseInt(latest_time, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	//调用service层的feed方法，传lastTime ,uidInt(这个用于判断isFavorite,仅对登录用户有效)
	videolist, err := service.Feed(ctx, c, lastTime, uidInt)
	//查询出错，范围code=1,
	if err != nil {
		hlog.Error("视频查询错误")
		c.JSON(consts.StatusInternalServerError, feedRes{
			Response:   Response{StatusCode: 1, StatusMsg: "查询视频错误"},
			Video_list: nil,
			Next_time:  nil,
		})
		return
	}
	var videoRes = []Video{}
	//	begin := time.Now().Unix()
	//将videolist数据copy到videoRes
	copyVideo(&videolist, &videoRes)
	//	end := time.Now().Unix()
	//hlog.Infof("copyVideo方法耗时:%v\n", end-begin)
	var next_time int64
	//取出next_time时间戳，因为是倒序，所以取切片的最后一个元素即可。
	if len(videoRes) != 0 {
		next_time = videoRes[len(videoRes)-1].Create_time.Unix()
	} else {
		next_time = time.Now().Unix()
	}
	//返回成功success代码
	c.JSON(consts.StatusOK, feedRes{
		Response:   Response{StatusCode: 0, StatusMsg: "success"},
		Video_list: videoRes,
		Next_time:  next_time,
	})
}

// 将service数据复刻到controller封装类  v1-->v2,指针操作
func copyVideo(v1 *[]service.Video, v2 *[]Video) {
	for _, temp := range *v1 {
		user := User{
			Id:            temp.UserID,
			Name:          temp.Username,
			FollowCount:   temp.FollowCount,
			FollowerCount: temp.FollowerCount,
			IsFollow:      temp.Is_follow,
		}
		var isFavorite bool
		if temp.IsFavorite != 0 {
			isFavorite = true
		}
		video := Video{
			Id:            temp.VideoID,
			Author:        user,
			PlayUrl:       temp.PlayURL,
			CoverUrl:      temp.CoverURL,
			FavoriteCount: temp.FavoriteCount,
			CommentCount:  temp.CommentCount,
			IsFavorite:    isFavorite,
			Title:         temp.Title,
			Create_time:   temp.CreateTime,
		}
		*v2 = append(*v2, video)
	}
}
