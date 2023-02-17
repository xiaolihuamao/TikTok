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

type PublishController struct {
	Response
	Video_list []Video `json:"video_list"`
}

// douyin/publish/action/ post
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
	files := form.File["data"]

	//title := c.Query("title")
	//token := c.Query("token")

	// 错误写法: data := os.File{"data"}
	//===============================================================
	//form, err := c.MultipartForm()
	//if err != nil {
	//	hlog.Error("文件格式有误")
	//	c.JSON(consts.StatusBadRequest, Response{
	//		StatusCode: -1,
	//		StatusMsg:  "请求参数错误",
	//	})
	//	return
	//}
	//==============================================================

	if token == nil || title == nil || len(token) != 1 || len(title) != 1 || token[0] == "" || title[0] == "" {
		c.JSON(consts.StatusNotFound, Response{
			StatusCode: -1,
			StatusMsg:  "(token ——> null) or (title ——> null)",
		})
		return
	}

	var uid interface{} //接收token解析出的uid
	//token不为空，解析token
	if token != nil {
		claims, err := mw.AuthMiddleware.GetClaimsFromJWT(ctx, c) //解析token,取出claims map
		if err != nil {
			hlog.Error("token解析错误，请使用正确的token")
			c.JSON(consts.StatusBadRequest, PublishController{
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
			uid = float64(-1)
			hlog.Error("token解析错误，请使用正确的token")
			c.JSON(consts.StatusBadRequest, PublishController{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "token解析错误，请使用正确的token",
				},
				Video_list: nil,
			})
			return
		}
		//token为空，给uid赋值float(-1)
	} else {
		uid = float64(-1)
		hlog.Info("未登录用户")
	}

	uidf := uid.(float64)
	uidInt := int64(uidf) //uid interface{}/float64--->int64 方便传参

	//for _, file := range files {
	//	fmt.Println(file.Filename)
	//	err := c.SaveUploadedFile(file, fmt.Sprintf("./file/upload/%s", file.Filename))
	//	if err != nil {
	//		hlog.Error("路径错误，上传失败")
	//		c.JSON(consts.StatusInternalServerError, Response{
	//			StatusCode: -1,
	//			StatusMsg:  "路径错误",
	//		})
	//		return
	//	}
	//}
	//c.JSON(consts.StatusOK, utils.H{
	//	"publish": "success",
	//})

	err = service.Publish(uidInt, files[0], c, title[0], ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, PublishController{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
			Video_list: nil,
		})
		return
	}
	//user := service.GetUserInfo(uidInt)

	c.JSON(consts.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})

}

// /douyin/publish/list/
func PublishList(ctx context.Context, c *app.RequestContext) {

	token := c.Query("token")
	userid := c.Query("user_id")

	//title := c.Query("title")

	//var uid interface{} //接收token解析出的uid
	//token和userid不为空，解析token
	//if token != "" {
	//	claims, err := mw.AuthMiddleware.GetClaimsFromJWT(ctx, c) //解析token,取出claims map
	//	if err != nil {
	//		hlog.Error("token解析错误，请使用正确的token")
	//	}
	//	//取出登录后返回的token中保存的uid---(interface{}/float64)
	//	uid = claims["id"]
	//	//若取不出uid,说明token错误或过期，给uid赋值float64(-1)
	//	if uid == nil {
	//		uid = float64(-1)
	//		hlog.Error("token解析错误，请使用正确的token")
	//	}
	//	//token为空，给uid赋值float(-1)
	//} else {
	//	uid = float64(-1)
	//	hlog.Info("未登录用户")
	//}
	//uidf := uid.(float64)
	//uidInt := int64(uidf) //uid interface{}/float64--->int64 方便传参

	if token == "" || userid == "" {
		c.JSON(consts.StatusNotFound, PublishController{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "(token ——> null) or (user_id ——> null)"},
		})
		return
	}

	//	err := service.PublishList(uidInt, user_id, title)
	//	if err != nil {
	//		hlog.Error("响应错误")
	//	}
	//	return
	//}
	//	user := service.GetUserInfo(uidInt)
	userId, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		hlog.Error("userid不合法")
		c.JSON(consts.StatusBadRequest, PublishController{
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

	//_, err = service.GetSnapshot("./test.mp4", "test", 1)
	//if err != nil {
	//	return
	//}

	c.JSON(consts.StatusOK, PublishController{
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
