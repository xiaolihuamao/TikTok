package controller

//author:zhuqitao
import (
	"TikTok/biz/model"
	mw "TikTok/biz/mw/jwt"
	redisUtil "TikTok/biz/mw/redis"
	"TikTok/biz/service"
	"TikTok/conf"
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

type UserResponse struct {
	Response
	User User `json:"user"`
}

type RegisterResponse struct {
	Response
	user_id int64  `json:"user_id"`
	token   string `json:"token"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || password == "" {
		hlog.Info("username或者password不能为空")
		c.JSON(consts.StatusInternalServerError, RegisterResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或者密码为空"},
			user_id:  0, //int返回默认值，未知
			token:    "",
		})
		return
	}
	_, err := service.Registeruser(ctx, c, username, password)
	if err != nil {
		hlog.Info("用户名重复或数据插入错误")
		c.JSON(consts.StatusInternalServerError, RegisterResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名重复或数据插入错误"},
			user_id:  0, //int返回默认值，未知
			token:    "",
		})
		return
	}
	mw.AuthMiddleware.LoginHandler(ctx, c)
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	id := c.GetInt64("user_id")
	if id == 0 {
		user_id := c.Query("user_id")
		id, _ = strconv.ParseInt(user_id, 10, 64)
	}

	//获取token
	token := c.Query("token")
	if token == "" {
		token = c.GetString("token")
	}
	if token == "" {
		hlog.Info("token 为空")
		c.JSON(consts.StatusInternalServerError, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户未登录"},
			User:     User{},
		})
		return
	}
	userlist, err := service.GetuserInfo(ctx, c, id)
	var userss = []User{}
	copyUser(&userlist, &userss)
	if err == nil && len(userss) != 0 {
		c.JSON(consts.StatusOK, UserResponse{
			Response: Response{StatusCode: 0, StatusMsg: "success"},
			User:     userss[0],
		})
	}

}
func copyUser(r1 *[]model.User, r2 *[]User) {
	for _, temp := range *r1 {
		followuser := User{
			Id:               temp.UserID,
			Name:             temp.Username,
			FollowCount:      temp.FollowCount,
			FollowerCount:    temp.FollowerCount,
			IsFollow:         true, //默认返回
			Avatar:           conf.IPAndPort + "/upload/backgrounds/20230219221523.jpg",
			Background_image: conf.IPAndPort + "/upload/backgrounds/20230219221607.jpg",
			Signature:        "曼曼女士的小木屋",
		}
		createUser(&followuser, followuser.Id)
		if followuser.Total_favorited == "" {
			followuser.Total_favorited = "0"
		}
		*r2 = append(*r2, followuser)
	}
}

func createUser(u *User, id int64) {
	pipe := redisUtil.Rdb
	userHashKey := fmt.Sprintf("userinfo_hash_%d", id)
	if pipe.Exists(userHashKey).Val() == 0 {
		return
	}
	u.Total_favorited = pipe.HGet(userHashKey, "Total_favorited").Val()
	Work_count, _ := pipe.HGet(userHashKey, "Work_count").Int64()
	u.Work_count = Work_count
	Favorite_count, _ := pipe.HGet(userHashKey, "Favorite_count").Int64()
	u.Favorite_count = Favorite_count
}
