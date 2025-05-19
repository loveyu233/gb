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


func TestDemo(t *testing.T) {
	// 创建配置
	config := Config{
		Title:       "API文档标题",
		Description: "API文档描述",
		Version:     "1.0.0",
		Host:        "localhost:8080",
		BasePath:    "/api/v1",
		OutputPath:  "./swagger.json",
	}

	// 创建生成器
	generator := NewGenerator(config)

	generator.AddAPI(APIInfo{
		Path:        "/auth/login",   // API路径
		Method:      "POST",          // 请求方法
		Summary:     "用户登录",          // API摘要
		Description: "用户登录接口",        // API描述
		Tags:        []string{"认证"},  // API标签
		Request:     LoginRequest{},  // 请求参数结构体
		Response:    LoginResponse{}, // 响应参数结构体
	})

	generator.AddAPI(APIInfo{
		Path:        "/auth/login/{id}",  // API路径
		Method:      "get",               // 请求方法
		Summary:     "url1",              // API摘要
		Description: "url1",              // API描述
		Tags:        []string{"url参数测试"}, // API标签
		PathParams: []ParamDescription{
			{
				Name:        "id",
				Description: "这是用户id",
				Type:        ParamTypeInteger,
			},
		},
		Response: LoginResponse{}, // 响应参数结构体
	})

	// 生成文档
	err := generator.Generate()
	if err != nil {
		fmt.Printf("生成Swagger文档失败: %v\n", err)
		return
	}
}

```