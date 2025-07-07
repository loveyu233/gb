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
	test2Routes := r.Group("/test2", gb.GinLogSetSkipLogFlag(), gb.GinLogSetModuleName("这是测试2模块"))
	{
		test2Routes.GET("/hello", gb.GinLogSetOptionName("hello"), func(c *gin.Context) {
			type Req struct {
				Age int64 `json:"age" form:"age" binding:"oneof=10 11 12"`
			}
			var req Req
			if err := c.BindQuery(&req); err != nil {
				gb.ResponseParamError(c, err)
				return
			}
			page, size := gb.ParsePaginationParams(c)
			gb.ResponseSuccess(c, fmt.Sprintf("hello %d %d", page, size))
		})
		test2Routes.GET("/page", func(c *gin.Context) {
			id, err := gb.ParseFromQuery(c, "id", "", gb.ParserInt64)
			if err != nil {
				gb.ResponseParamError(c, err)
				return
			}
			gb.ResponseSuccess(c, id)
		})
		test2Routes.GET("/path/:id", func(c *gin.Context) {
			id, err := gb.ParseFromPath(c, "id1", gb.ParserInt64)
			if err != nil {
				gb.ResponseParamError(c, err)
				return
			}
			gb.ResponseSuccess(c, id)
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
		TokenService: gb.NewJWTTokenService[TokenTestUser]("adadasdasdasdasdasd", gb.WithRedisClient[TokenTestUser](gb.RedisClient), gb.WithRedisTokenCheck[TokenTestUser](func(token string) string {
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
	)
}
