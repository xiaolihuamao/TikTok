package controller

//@author:zhangshuo
import (
	"TikTok/biz/model"
	"TikTok/biz/service"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CommentActionResponse
// 进行评论操作的返回结构体
type CommentActionResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	Comment    CommentInfo `json:"comment"`
}

// CommentListResponse
// 获取评论列表的返回结构体
type CommentListResponse struct {
	StatusCode  int32         `json:"status_code"`
	StatusMsg   string        `json:"status_msg,omitempty"`
	CommentList []CommentInfo `json:"comment_list,omitempty"`
}

// CommentAction
// 评论操作函数
func CommentAction(ctx context.Context, c *app.RequestContext) {
	log.Printf("the CommentAction function is running") //提示函数正在运行

	//获取userId
	userid := c.Query("user_id")
	userId, err := strconv.ParseInt(userid, 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "comment userId json invalid",
		})
		log.Println("CommentController-Comment_Action: return comment userId json invalid") //函数返回userId无效
		return
	}
	log.Printf("userId:%v", userId)

	//获取videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "comment userId json invalid",
		})
		log.Println("CommentController-Comment_Action: return comment videoId json invalid")
		return
	}
	log.Printf("videoId:%v", videoId)

	//获取操作类型
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	//错误处理
	if err != nil || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "comment actionType json invalid",
		})
		log.Println("CommentController-Comment_Action: return actionType json invalid") //评论类型数据无效
		return
	}
	log.Printf("actionType:%v", actionType)

	//调用service层评论函数，完成发送或删除评论
	if actionType == 1 { //actionType为1，则进行发表评论操作
		content := c.Query("comment_text")
		//发表评论数据准备
		var sendComment model.Comment
		sendComment.UserID = userId
		sendComment.VideoID = videoId
		sendComment.Content = content
		timeNow := time.Now()
		// 格式有问题，暂置
		sendComment.CreateDate = timeNow.String()
		//发表评论
		commentInfo, err := service.Send(sendComment)
		//发表评论失败
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "send comment failed",
			})
			log.Println("CommentController-Comment_Action: return send comment failed") //发表失败
			return
		}
		//发表评论成功:
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "send comment success",
			Comment:    commentInfo,
		})
		log.Println("CommentController-Comment_Action: return Send success") //发表评论成功，返回正确信息
		return
	} else {
		//actionType为2，则进行删除评论操作
		//获取要删除的评论的id
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "delete commentId invalid",
			})
			log.Println("CommentController-Comment_Action: return commentId invalid") //评论id格式错误
			return
		}
		log.Printf("commentId:%v", commentId)

		//删除评论操作
		err = service.DelComment(commentId)
		if err != nil { //删除评论失败
			str := err.Error()
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  str,
			})
			log.Println("CommentController-Comment_Action: return delete comment failed") //删除失败
			return
		}
		//删除评论成功
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "delete comment success",
		})

		log.Println("CommentController-Comment_Action: return delete success") //函数执行成功，返回正确信息
		return
	}

}

func CommentList(ctx context.Context, c *app.RequestContext) {
	log.Println("CommentController-Comment_List: running") //函数已运行
	//获取userId
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "comment userId json invalid",
		})
		log.Println("CommentController-Comment_Action: return comment userId json invalid") //userId无效
		return
	}
	log.Printf("userId:%v", userId)

	//获取videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "comment videoId json invalid",
		})
		log.Println("CommentController-Comment_List: return videoId json invalid") //视频id格式有误
		return
	}
	log.Printf("videoId:%v", videoId)

	//调用service层评论函数
	commentList, err := service.GetList(videoId, userId)
	if err != nil { //获取评论列表失败
		c.JSON(http.StatusOK, CommentListResponse{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		log.Println("CommentController-Comment_List: return list false") //查询列表失败
		return
	}

	//获取评论列表成功
	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "get comment list success",
		CommentList: commentList,
	})
	log.Println("CommentController-Comment_List: return success") //成功返回列表
	return
}
