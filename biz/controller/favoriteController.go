package controller

//@author:zhangxiyang
import (
	mw "TikTok/biz/mw/jwt"
	"TikTok/biz/service"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	//参数校验
	if token == "" || videoId == "" || (actionType != "1" && actionType != "2") {
		c.JSON(consts.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  "参数缺少",
		})
		return
	}
	vid, _ := strconv.ParseInt(videoId, 0, 64)
	var uid interface{}
	claims, err := mw.AuthMiddleware.GetClaimsFromJWT(ctx, c) //解析token,取出claims map
	if err != nil {
		hlog.Error("token解析错误，请使用正确的token")
		c.JSON(consts.StatusBadRequest, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "token解析错误，请使用正确的token",
			},
			Video_list: nil,
		})
		return
	}
	//取出登录后返回的token中保存的uid---(interface{}/float64)
	uid = claims["id"]
	//若取不出uid,说明token错误或过期，给uid赋值float64(-1)
	if uid == nil {
		hlog.Error("token解析错误，请使用正确的token")
		c.JSON(consts.StatusBadRequest, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "token解析错误，请使用正确的token",
			},
			Video_list: nil,
		})
		return
	}
	uidf := uid.(float64)
	uidInt := int64(uidf) //uid interface{}/float64--->int64 方便传参
	var errS error
	if actionType == "1" {
		errS = service.AddLike(uidInt, vid, ctx)
	} else {
		errS = service.CancelLike(uidInt, vid, ctx)
	}
	if errS != nil {
		c.JSON(consts.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "服务器错误",
		})
		return
	}
	c.JSON(consts.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

func FavoriteList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	if user_id == "" || token == "" {
		c.JSON(consts.StatusBadRequest, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "请求错误",
			},
			Video_list: nil,
		})
		return
	}
	fmt.Println(token)
	uidInt, _ := strconv.ParseInt(user_id, 0, 64)
	//调用service查询
	videoInfo := service.FavorList(uidInt, ctx)
	var videoList []Video
	copyFavorList(videoInfo, &videoList)
	c.JSON(consts.StatusOK, publishRes{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Video_list: videoList,
	})
}

func copyFavorList(videoInfo []service.Video, videoList *[]Video) {
	for _, temp := range videoInfo {
		author := User{
			Id:            temp.UserID,
			Name:          temp.Username,
			FollowCount:   temp.FollowCount,
			FollowerCount: temp.FollowerCount,
			IsFollow:      true,
		}
		video := Video{
			Id:            temp.VideoID,
			Author:        author,
			PlayUrl:       temp.PlayURL,
			CoverUrl:      temp.CoverURL,
			FavoriteCount: temp.FavoriteCount,
			CommentCount:  temp.CommentCount,
			IsFavorite:    true,
			Title:         temp.Title,
			Create_time:   temp.CreateTime,
		}
		*videoList = append(*videoList, video)
	}
}
