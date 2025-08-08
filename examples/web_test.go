package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
)

func init() {
	gb.PublicRoutes = append(gb.PublicRoutes, registerDemo1PublicRoutes)
	gb.PrivateRoutes = append(gb.PrivateRoutes, registerDemo1PrivateRoutes)
}

type Ctx struct {
	Username string   `json:"username"`
	Age      int64    `json:"age"`
	Address  []string `json:"address"`
}

func registerDemo1PrivateRoutes(r *gin.RouterGroup) {
	testRoutes := r.Group("/test1", gb.GinLogSetModuleName("这是测试1模块"))
	{
		testRoutes.POST("/world", gb.GinLogSetOptionName("world"), func(c *gin.Context) {
			// 需要启动  TestToken 测试可以查看
			//tu, ex := gb.GetGinContextValue[TokenTestUser](c, "data")
			data, _ := gb.GetGinContextTokenLoadData[TokenTestUser](c)
			tokenString := gb.GetGinContextTokenString(c)
			gb.ResponseSuccess(c, map[string]any{
				"data":  data,
				"token": tokenString,
			})
		})
	}

}

func registerDemo1PublicRoutes(r *gin.RouterGroup) {
	//  gb.GinLogSetSkipLogFlag(), 不输出请求信息
	test2Routes := r.Group("/test2", gb.GinLogSetModuleName("这是测试2模块"), func(c *gin.Context) {
		c.Set("id", Ctx{
			Username: "username-1",
			Age:      11,
			Address:  []string{"a", "b", "c"},
		})
	})
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
		test2Routes.GET("/page/:id", func(c *gin.Context) {
			id, err := gb.GetGinQueryDefault[int64](c, "id", 1)
			if err != nil {
				gb.ResponseParamError(c, err)
				return
			}
			pathID, err := gb.GetGinPathRequired[int64](c, "id")
			if err != nil {
				gb.ResponseParamError(c, err)
				return
			}
			gb.ResponseSuccess(c, map[string]any{
				"query": id,
				"path":  pathID,
			})
		})
		test2Routes.GET("/path/:id", func(c *gin.Context) {
			id, exists := gb.GetGinContextValue[Ctx](c, "id")
			gb.ResponseSuccess(c, map[string]any{
				"value":  id,
				"exists": exists,
			})
		})
	}
}

type TokenTestUser struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

var authConfig *gb.GinAuthConfig

func TestZiDingYiToken(t *testing.T) {
	gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}))

	authConfig = &gb.GinAuthConfig{
		TokenService: gb.NewJWTTokenService("adadasdasdasdasdasd", gb.WithRedisTokenKey(func(token string) string {
			return fmt.Sprintf("zidingyikey:%s", token)
		})),
	}

	t.Log(authConfig.TokenService.Generate(TokenTestUser{ID: 1, Username: "hzyyy"}, 1000*time.Second))

	gb.InitHTTPServerAndStart(authConfig, "127.0.0.1:8080")
}

func TestPublicHttp(t *testing.T) {
	gb.InitPublicHTTPServerAndStart("127.0.0.1:8888")
}
