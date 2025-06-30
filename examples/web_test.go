package examples

import (
	"errors"
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

type User struct {
	ID   int
	Name string
}

func TestHttp(t *testing.T) {
	u := &User{ID: 1, Name: "test"}
	t.Log(gb.NewJWTTokenService("").Generate(u, 1000*time.Second))
	gb.InitHTTPServerAndStart(":8080", gb.WithGinRouterModel("release"))
}

func TestParam(t *testing.T) {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		type Req struct {
			Username string `json:"username" binding:"required"`
		}
		var req Req
		if err := c.ShouldBindQuery(&req); err != nil {
			gb.ResponseParamError(c, err)
		}
	})
	engine.Run("127.0.0.1:8080")
}

func TestErr(t *testing.T) {
	appError := gb.ConvertToAppError(errors.New("test error"))
	fmt.Println(appError)
}

func TestLog(t *testing.T) {
	engine := gin.Default()
	engine.Use(gb.MiddlewareLogger(gb.MiddlewareLogConfig{HeaderKeys: []string{"token"}}))

	engine.GET("/a", func(c *gin.Context) {
		logger := gb.GetContextLogger(c)
		logger.Info().Str("adadadasdasdas", "adasdasdasdasdasdas").Msg("adasdasdasdasdasdas")
		gb.ResponseSuccess(c, map[string]any{
			"test": "test",
		})
	})

	engine.Run("127.0.0.1:8080")
}

type TokenTestUser struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

func TestToken(t *testing.T) {
	t.Log(gb.DefaultGinTokenConfig.TokenService.Generate(&TokenTestUser{Username: "hzyy", ID: 19}, 1000*time.Second))
	gb.InitHTTPServerAndStart("127.0.0.1:8080",
		gb.WithGinRouterModel(gb.GinModelDebug),
		gb.WithGinRouterSkipHealthzLog(),
		gb.WithGinRouterSkipApiMap("/api/test2/hello"),
		gb.WithGinRouterTokenData(new(TokenTestUser)),
	)
}

func TestZiDingYiToken(t *testing.T) {
	// 不初始化就会使用默认的
	gb.InitCustomGinAuthConfig(&gb.GinAuthConfig{
		DataPtr:      new(TokenTestUser),
		TokenService: gb.NewJWTTokenService("adadasdasdasdasdasd"),
		GetTokenStrFunc: func(c *gin.Context) string {
			return c.GetHeader("jwt-token")
		},
		HandleError: gb.DefaultGInTokenErrHandler,
	})

	// Generate的值必须和gb.InitCustomGinAuthConfig的DataPtr是一样的
	t.Log(gb.CustomGinAuthConfig.TokenService.Generate(&TokenTestUser{Username: "hzyy", ID: 19}, 1000*time.Second))
	gb.InitHTTPServerAndStart("127.0.0.1:8080",
		gb.WithGinRouterModel(gb.GinModelDebug),
		gb.WithGinRouterSkipHealthzLog(),
		gb.WithGinRouterSkipApiMap("/api/test2/hello"),
	)
}
