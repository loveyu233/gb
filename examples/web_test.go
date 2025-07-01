package examples

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
	"testing"
	"time"
)

func init() {
	gb.PublicRoutes = append(gb.PublicRoutes, registerDemo1PublicRoutes)
	gb.PrivateRoutes = append(gb.PrivateRoutes, registerDemo1PrivateRoutes)
}

func registerDemo1PrivateRoutes(r *gin.RouterGroup) {
	testRoutes := r.Group("/test1", gb.GinLogSetModuleName("这是测试1模块"))
	{
		testRoutes.POST("/world", gb.GinLogSetOptionName("world"), func(c *gin.Context) {
			// 需要启动  TestToken 测试可以查看
			value, exists := c.Get("data")
			if exists {
				if user, ok := value.(*TokenTestUser); ok {
					// 测试删除token
					err := authConfig.TokenService.DeleteRedisToken(c.GetHeader("jwt-token"))
					fmt.Println(err)
					gb.ResponseSuccess(c, user)
					return
				}
			}

			gb.ResponseSuccess(c, "world")
		})
	}

}

func registerDemo1PublicRoutes(r *gin.RouterGroup) {
	test2Routes := r.Group("/test2", gb.GinLogSetModuleName("这是测试2模块"))
	{
		test2Routes.GET("/hello", gb.GinLogSetOptionName("hello"), func(c *gin.Context) {
			page, size := gb.ParsePaginationParams(c, map[string]int{"page": 0, "size": 20})
			gb.ResponseSuccess(c, fmt.Sprintf("hello %d %d", page, size))
		})
	}
}

type TokenTestUser struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

var authConfig *gb.GinAuthConfig[TokenTestUser]

func TestZiDingYiToken(t *testing.T) {
	gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}))

	authConfig = &gb.GinAuthConfig[TokenTestUser]{
		DataPtr: new(TokenTestUser),
		TokenService: gb.NewJWTTokenService[TokenTestUser]("adadasdasdasdasdasd", gb.WithRedisClient[TokenTestUser](gb.RedisClient), gb.WithRedisTokenCheck[TokenTestUser](true, func(token string) string {
			return fmt.Sprintf("zidingyikey:%s", token)
		})),
		GetTokenStrFunc: func(c *gin.Context) string {
			return c.GetHeader("jwt-token")
		},
		HandleError: gb.DefaultGInTokenErrHandler,
	}

	t.Log(authConfig.TokenService.Generate(TokenTestUser{ID: 1, Username: "hzyyy"}, 1000*time.Second))

	gb.InitHTTPServerAndStart(authConfig, "127.0.0.1:8080",
		gb.WithGinRouterModel(gb.GinModelDebug),
		gb.WithGinRouterSkipHealthzLog(),
		gb.WithGinRouterSkipApiMap("/api/test2/hello"),
	)
}
