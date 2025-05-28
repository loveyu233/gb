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
	testRoutes := r.Group("/test1", gb.SetModuleName("这是测试1模块"))
	{
		testRoutes.POST("/world", gb.SetOptionName("world"), func(c *gin.Context) {
			gb.ResponseSuccess(c, "world")
		})
	}

}

func registerDemo1PublicRoutes(r *gin.RouterGroup) {
	test2Routes := r.Group("/test2", gb.SetModuleName("这是测试2模块"))
	{
		test2Routes.GET("/hello", gb.SetOptionName("hello"), func(c *gin.Context) {
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
	gb.InitRouter("debug", "/abc", []gin.HandlerFunc{gb.GinAuth(&User{}, gb.DefaultGinConfig)}, gb.AddTraceID(), gb.AddRequestTime(), gb.ResponseLogger(), gb.GinRecovery(true))
	server := gb.CreateHTTPServer("127.0.0.1:8080")
	go gb.StartHTTPServer(server)
	gb.SetupGracefulShutdown(server)
}
