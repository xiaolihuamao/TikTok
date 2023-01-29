package mw

//@author fuxingyuan
import (
	"TikTok/biz/controller"
	"TikTok/biz/dao"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"log"
	"time"
)

type login struct {
	Username string `query:"username,required" json:"username,required"`
	Password string `query:"password,required" json:"password,required"`
}
type loginRes struct {
	controller.Response
	Token   interface{} `form:"token" json:"token"`
	User_id interface{} `form:"user_id" json:"user_Id"`
}

var identityKey = "id"

type UserDemo struct {
	UserName  string
	FirstName string
	LastName  string
}

var AuthMiddleware *jwt.HertzJWTMiddleware
var errjwt error

func Initjwt() {
	// the jwt middleware
	AuthMiddleware, errjwt = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Minute,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserDemo); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &UserDemo{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginVals login
			if err := c.BindAndValidate(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			user := dao.Use(dao.Db).User
			validUser, err := user.WithContext(ctx).Where(user.Username.Eq(userID), user.Password.Eq(password)).First()
			if err == nil && validUser != nil {
				c.Set("id", validUser.UserID)
				return &UserDemo{
					UserName:  userID,
					LastName:  "Hertz",
					FirstName: "CloudWeGo",
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, loginRes{
				User_id: nil,
				Token:   nil,
				Response: controller.Response{
					StatusCode: -1,
					StatusMsg:  message,
				}})
		},
		TokenLookup: "query:token,form:token,param:token",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
			userid, _ := c.Get("id")
			c.JSON(code, loginRes{
				User_id: userid,
				Token:   message,
				Response: controller.Response{
					StatusCode: 0,
					StatusMsg:  "success login!",
				}})
		},
	})
	if errjwt != nil {
		log.Fatal("JWT Error:" + errjwt.Error())
	}
	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := AuthMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
}
