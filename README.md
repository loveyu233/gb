```go
package gb

import (
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func init() {
	PublicRoutes = append(PublicRoutes, registerDemo1PublicRoutes)
	PrivateRoutes = append(PrivateRoutes, registerDemo1PrivateRoutes)
}

func registerDemo1PrivateRoutes(r *gin.RouterGroup) {
	testRoutes := r.Group("/test1", SetModuleName("这是测试1模块"))
	{
		testRoutes.POST("/world", SetOptionName("world"), func(c *gin.Context) {
			ResponseSuccess(c, "world")
		})
	}

}

func registerDemo1PublicRoutes(r *gin.RouterGroup) {
	test2Routes := r.Group("/test2", SetModuleName("这是测试2模块"))
	{
		test2Routes.GET("/hello", SetOptionName("hello"), func(c *gin.Context) {
			ResponseSuccess(c, "hello")
		})
	}
}

type User struct {
	ID   int
	Name string
}

func TestHttp(t *testing.T) {
	u := &User{ID: 1, Name: "test"}
	t.Log(NewJWTTokenService("").Generate(u, 1000*time.Second))
	InitRouter("debug", "/abc", 
		[]gin.HandlerFunc{GinAuth(&User{}, defaultGinConfig)}, 
		AddTraceID(), AddRequestTime(), ResponseLogger(), GinRecovery(true), 
    )
	server := CreateHTTPServer("127.0.0.1:8080")
	go StartHTTPServer(server)
	SetupGracefulShutdown(server)
}

```