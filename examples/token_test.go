package examples

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
)

func TestT1(t *testing.T) {
	var identity = "user_id"
	middleware, err := gb.InitGinJWTMiddleware(&gb.GinJWTMiddleware{
		Key:         []byte("abc123abc123....."),
		TokenLookup: "header:Authorization",
		Timeout:     24 * time.Hour,
		IdentityKey: identity,
		// 登录时验证密码
		Authenticator: func(c *gin.Context) (interface{}, error) {
			t.Log("Authenticator")
			if c.Query("username") == "root" {
				return 1001, nil
			}
			return nil, gb.ErrBadRequest.WithMessage("账号错误")
		},
		// 在令牌中设置负载
		PayloadFunc: func(data interface{}) gb.MapClaims {
			t.Log("PayloadFunc", data)
			return gb.MapClaims{
				identity: data,
			}
		},
		// 如果验证密码成功，则生成登录响应
		LoginResponse: func(c *gin.Context, code int, token string, time time.Time) {
			t.Log("LoginResponse", code, token, time.String())
			c.Set("token", token)
		},
		// 验证令牌并获取登录用户的id
		Authorizator: func(data interface{}, c *gin.Context) bool {
			t.Log("Authorizator", data)
			return true
		},
		// 验证失败，生成消息
		Unauthorized: func(c *gin.Context, code int, message string) {
			t.Log("Unauthorized", code, message)
			gb.ResponseError(c, gb.ErrTokenInvalid)
		},
		// JWT中间件失败时的HTTP状态消息
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			t.Log(e)
			return gb.ErrTokenInvalid.Message
		},
	})
	if err != nil {
		t.Log(err)
		return
	}
	gb.PublicRoutes = append(gb.PrivateRoutes, func(group *gin.RouterGroup) {
		group.GET("/p1", func(c *gin.Context) {
			gb.ResponseSuccess(c, "p1")
		})
		group.POST("/login", middleware.LoginHandler(), func(c *gin.Context) {
			gb.ResponseSuccess(c, c.GetString("token"))
		})
	})
	gb.PrivateRoutes = append(gb.PrivateRoutes, func(group *gin.RouterGroup) {
		group.GET("/p2", func(c *gin.Context) {
			gb.ResponseSuccess(c, c.GetString("token"))
		})
	})
	gb.InitHTTPServerAndStart("localhost:8888", gb.WithGinRouterAuthHandler(middleware.MiddlewareFunc()))
}
