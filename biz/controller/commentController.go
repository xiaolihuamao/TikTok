package controller

//@author:zhangshuo
import (
	"TikTok/biz/model"
	mw "TikTok/biz/mw/jwt"
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

type CommentRes struct {
	Response
	service.Comment
}
type CommentListRes struct {
	Response
	CommentList []service.Comment `json:"comment_list"`
}

// /douyin/comment/action/ post
// query ?name=value&name=value
func CommentAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	if action_type == "" || token == "" || video_id == "" {
		c.JSON(consts.StatusNotFound, CommentRes{
			Response: Response{StatusCode: -1, StatusMsg: "参数不足"},
			Comment:  service.Comment{},
		})
		return
	}
	//==========================
	var uid interface{}                                       //接收token解析出的uid
	claims, err := mw.AuthMiddleware.GetClaimsFromJWT(ctx, c) //解析token,取出claims map
	if err != nil {
		hlog.Error("token解析错误，请使用正确的token")
	}
	//取出登录后返回的token中保存的uid---(interface{}/float64)
	uid = claims["id"]
	if uid == nil {
		c.JSON(consts.StatusBadRequest, CommentRes{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "error token",
			},
			Comment: service.Comment{},
		})
		return
	}
	uidf := uid.(float64)
	uidInt := int64(uidf)
	//一切参数就绪，调用service方法
	videoid, _ := strconv.Atoi(video_id)
	//增加评论
	if action_type == "1" {
		comment_text := c.Query("comment_text")
		if comment_text == "" {
			hlog.Error("评论内容为空")
			c.JSON(consts.StatusNotFound, CommentRes{
				Response: Response{StatusCode: -1, StatusMsg: "评论为空"},
				Comment:  service.Comment{},
			})
			return
		}

		err := service.AddComment(int64(videoid), comment_text, uidInt)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, CommentRes{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  "评论失败",
				},
				Comment: service.Comment{},
			})
			return
		}
		c.JSON(consts.StatusOK, CommentRes{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			Comment: service.Comment{
				Comment: model.Comment{
					Content: comment_text,
				},
			},
		})
		//删除评论
	} else if action_type == "2" {
		comment_id := c.Query("comment_id")
		if comment_id == "" {
			c.JSON(consts.StatusBadRequest, CommentRes{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  "请求参数缺失",
				},
				Comment: service.Comment{},
			})
			return
		}
		commentId, _ := strconv.Atoi(comment_id)
		err := service.DelComment(int64(commentId), int64(videoid))
		if err != nil {
			hlog.Error("删除失败")
			c.JSON(consts.StatusInternalServerError, CommentRes{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  "删除评论失败",
				},
				Comment: service.Comment{},
			})
			return
		}
		c.JSON(consts.StatusOK, CommentRes{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			Comment: service.Comment{},
		})
	}
}

func CommentList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	if token == "" || video_id == "" {
		c.JSON(consts.StatusBadRequest, CommentRes{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "请求参数缺失",
			},
			Comment: service.Comment{},
		})
		return
	}
	//===========
	//解析token
	//===========
	videoId, _ := strconv.ParseInt(video_id, 0, 64)
	cList := service.CommentList(videoId)
	c.JSON(consts.StatusOK, CommentListRes{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		CommentList: cList,
	})
}
