package mw

//@author fuxingyuan
import (
	"TikTok/biz/dao"
	"TikTok/biz/service"
	"TikTok/conf"
	"context"
	"encoding/base64"
	"fmt"
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
	StatusCode int32       `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	Token      interface{} `form:"token" json:"token"`
	User_id    int         `form:"user_id" json:"user_id"`
}

var identityKey = "id"

type UserDemo struct {
	UserName  string
	FirstName string
	LastName  string
	Uid       int64
}

var AuthMiddleware *jwt.HertzJWTMiddleware
var errjwt error

func Initjwt() {
	// the jwt middleware
	AuthMiddleware, errjwt = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*UserDemo); ok {
				return jwt.MapClaims{
					identityKey: v.Uid,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			uidFloat := claims[identityKey].(float64)
			fmt.Println(uidFloat)
			uid := int(uidFloat)
			uid64 := int64(uid)
			return &UserDemo{
				Uid:       uid64,
				UserName:  "",
				FirstName: identityKey,
				LastName:  identityKey,
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
			validUser, err := user.WithContext(ctx).Select(user.UserID, user.Password).Where(user.Username.Eq(userID)).First()
			if err != nil || validUser == nil {
				return nil, jwt.ErrFailedAuthentication
			}
			key1 := []byte(conf.Secret)
			result1 := service.DesCbcEncryption([]byte(password), key1)
			passwordaes := base64.StdEncoding.EncodeToString(result1)
			if passwordaes != validUser.Password {
				return nil, jwt.ErrFailedAuthentication
			}
			c.Set("id", validUser.UserID)
			return &UserDemo{
				Uid: validUser.UserID,
			}, nil
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, loginRes{
				User_id:    -1,
				Token:      nil,
				StatusCode: -1,
				StatusMsg:  message,
			})
		},
		TokenLookup: "query:token,form:token,param:token",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, message string, time time.Time) {
			userid := c.GetInt64("id")
			c.JSON(code, loginRes{
				User_id:    int(userid),
				Token:      message,
				StatusCode: 0,
				StatusMsg:  "success login!",
			})
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
