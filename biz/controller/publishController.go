package controller

//author:zhangwangjun
import (
	mw "TikTok/biz/mw/jwt"
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

type publishRes struct {
	Response
	Video_list []Video `json:"video_list"`
}

// /douyin/publish/action
func Publish(ctx context.Context, c *app.RequestContext) {
	form, err := c.MultipartForm()
	if err != nil {
		hlog.Error("文件格式有误")
		c.JSON(consts.StatusBadRequest, Response{
			StatusCode: -1,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	token := form.Value["token"]

	title := form.Value["title"]

	if token == nil || len(token) != 1 || title == nil || len(title) != 1 || token[0] == "" || title[0] == "" {
		hlog.Error("无token或标题缺少")
		c.JSON(consts.StatusBadRequest, Response{
			StatusCode: -1,
			StatusMsg:  "无token或标题缺少",
		})
		return
	}
	//取出token中的uid
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
	files := form.File["data"]
	if files == nil || len(files) == 0 {
		hlog.Error("文件传入错误")
		c.JSON(consts.StatusBadRequest, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "文件传入错误",
			},
			Video_list: nil,
		})
		return
	}
	err = service.Publish(uidInt, files[0], c, title[0], ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			Video_list: nil,
		})
		return
	}
	//插入成功
	c.JSON(consts.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

func PublishList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	userid := c.Query("user_id")
	//参数校验
	if token == "" || userid == "" {
		c.JSON(consts.StatusBadRequest, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "请求错误",
			},
			Video_list: nil,
		})
		return
	}
	userId, err := strconv.ParseInt(userid, 0, 64)
	if err != nil {
		hlog.Error("userid不合法")
		c.JSON(consts.StatusBadRequest, publishRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "userid不合法",
			},
			Video_list: nil,
		})
		return
	}
	videoInfo := service.PublishList(userId, ctx)
	var videoList []Video
	if len(videoInfo) != 0 {
		copyPublishList(videoInfo, &videoList)
	}
	c.JSON(consts.StatusOK, publishRes{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		Video_list: videoList,
	})
}

// 将service的数据封装到controller层
func copyPublishList(videoInfo []service.Video, videoList *[]Video) {
	author := User{
		Id:            videoInfo[0].UserID,
		Name:          videoInfo[0].Username,
		FollowCount:   videoInfo[0].FollowCount,
		FollowerCount: videoInfo[0].FollowerCount,
		IsFollow:      true,
	}

	for _, temp := range videoInfo {
		var isFavorite bool
		if temp.IsFavorite != 0 {
			isFavorite = true
		}
		video := Video{
			Id:            temp.VideoID,
			Author:        author,
			PlayUrl:       temp.PlayURL,
			CoverUrl:      temp.CoverURL,
			FavoriteCount: temp.FavoriteCount,
			CommentCount:  temp.CommentCount,
			IsFavorite:    isFavorite,
			Title:         temp.Title,
			Create_time:   temp.CreateTime,
		}
		*videoList = append(*videoList, video)
	}
}
