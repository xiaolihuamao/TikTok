package service

//author:zhangshuo
import (
	"TikTok/biz/dao"
	"TikTok/conf"
	"log"
	"sort"
	"sync"
)

// 定义相关数据结构

type CommentInfo struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

/*
	评论模块自己request实现的方法：
*/

// Send
// 发表评论，传进来评论的基本信息，返回保存是否成功的状态描述
func Send(comment dao.CommentData) (CommentInfo, error) {
	log.Println("CommentService-Send: running") //函数已运行
	//数据准备
	var commentInfo dao.CommentData
	commentInfo.VideoId = comment.VideoId         //评论视频id传入
	commentInfo.UserId = comment.UserId           //评论用户id传入
	commentInfo.CommentText = comment.CommentText //评论内容传入
	commentInfo.CreateDate = comment.CreateDate   //评论时间
	commentInfo.Cancel = conf.ValidComment        //评论状态，0，有效

	//评论信息存储：
	commentRtn, err := dao.InsertComment(commentInfo)
	if err != nil {
		return CommentInfo{}, err
	}

	//查询用户信息
	userData, err2 := GetUserByIdWithCurId(comment.UserId, comment.UserId)
	if err2 != nil {
		return CommentInfo{}, err2
	}

	//拼接,
	commentData := CommentInfo{
		Id:         commentRtn.Id,
		User:       userData,
		Content:    commentRtn.CommentText,
		CreateDate: commentRtn.CreateDate.Format(conf.DateTime),
	}

	//返回结果
	return commentData, nil
}

// DelComment
// 删除评论，传入评论id即可
func DelComment(commentId int64) error {
	log.Println("CommentService-DelComment: running") //函数已运行
	//直接走数据库删除
	return dao.DeleteComment(commentId)
}

// GetList
// 查看评论列表-返回评论list
func GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	log.Println("CommentService-GetList: running") //函数已运行
	//调用dao，先查评论，再循环查用户信息：
	//先查询评论列表信息
	commentList, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
		return nil, err
	}
	//当前有0条评论
	if commentList == nil {
		return nil, nil
	}

	//提前定义好切片长度
	commentInfoList := make([]CommentInfo, len(commentList))

	wg := &sync.WaitGroup{}
	wg.Add(len(commentList))
	idx := 0
	for _, comment := range commentList {
		//调用方法组装评论信息，再append
		var commentData CommentInfo
		//将评论信息进行组装，添加想要的信息,插入从数据库中查到的数据
		go func(comment dao.CommentData) {
			oneComment(&commentData, &comment, userId)
			commentInfoList[idx] = commentData
			idx = idx + 1
			wg.Done()
		}(comment)
	}
	wg.Wait()
	//评论排序-按照主键排序
	sort.Sort(CommentSlice(commentInfoList))

	log.Println("CommentService-GetList: return list success") //函数执行成功，返回正确信息
	return commentInfoList, nil
}

// 此函数用于给一个评论赋值：评论信息+用户信息 填充
func oneComment(comment *CommentInfo, com *dao.CommentData, userId int64) {
	var wg sync.WaitGroup
	wg.Add(1)
	//根据评论用户id和当前用户id，查询评论用户信息
	var err error
	comment.Id = com.Id
	comment.Content = com.CommentText
	comment.CreateDate = com.CreateDate.Format(conf.DateTime)
	comment.User, err = GetUserByIdWithCurId(com.UserId, userId)
	if err != nil {
		log.Println("CommentService-GetList: GetUserByIdWithCurId return err: " + err.Error()) //函数返回提示错误信息
	}
	wg.Done()
	wg.Wait()
}

// GetUserByIdWithCurId
// 已登录情况下,根据user_id获得User对象
func GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	user := User{
		Id:            0,
		Name:          "",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	//通过id获取TableUser结构体
	tableUser, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")

	//没有社交功能接口，社交相关数据均返回0
	user = User{
		Id:            id,
		Name:          tableUser.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return user, nil
}

// CommentSlice 此变量以及以下三个函数都是做排序-准备工作
type CommentSlice []CommentInfo

func (a CommentSlice) Len() int { //重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { //重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { //重写Less()方法
	return a[i].Id > a[j].Id
}
