package examples

import (
	"testing"

	"github.com/loveyu233/gb"
)

func TestSwagger(t *testing.T) {
	swaggerGenerator := gb.NewSwaggerGenerator(gb.SwaggerGlobalConfig{
		Title:       "这是标题",
		Description: "这是描述",
		Version:     "v1.0.0",
		Host:        "127.0.0.1",
		BasePath:    "/api/v1",
		Schemes:     []string{"http", "https"},
		OutputPath:  "", //这是输出目录,不设置默认 swagger/swagger.json
	})
	// 添加全局参数
	swaggerGenerator.AddGlobalHeaderParams([]gb.SwaggerParamDescription{{
		Name:        "trace_id",
		Description: "请求链路追踪id",
		Type:        gb.ParamTypeString,
		Required:    false,
	}})
	type UserReq struct {
		Page int `json:"page" desc:"页码"`
		Size int `json:"size" desc:"页大小"`
	}
	type User struct {
		Username string `json:"username" desc:"用户名"`
		Age      int    `json:"age" desc:"年龄"`
	}
	type UserRes struct {
		Code  int    `json:"code" desc:"响应状态码"`
		Msg   string `json:"msg" desc:"描述"`
		Users []User `json:"users" desc:"返回的全部user数据"`
	}
	swaggerGenerator.AddAPI(gb.SwaggerAPIInfo{
		Path:        "/user",
		Method:      "get", //大小写不敏感
		Summary:     "获取用户列表",
		Description: "获取用户列表下的全部数据",
		Tags:        []string{"user"},
		Response:    UserRes{},
		QueryParams: UserReq{},
	})
	swaggerGenerator.AddAPI(gb.SwaggerAPIInfo{
		Path:        "/user/{id}",
		Method:      "get", //大小写不敏感
		Summary:     "获取指定用户信息",
		Description: "获取指定用户信息",
		Tags:        []string{"user"},
		Response:    User{},
		PathParams: []gb.SwaggerParamDescription{{
			Name:        "id",
			Description: "用户id",
			Type:        gb.ParamTypeInteger,
			Required:    true,
		}},
	})
	type CreateUser struct {
		Username string `json:"username" binding:"required" desc:"用户名"`
		Age      int    `json:"age" desc:"年龄"`
		Password string `json:"password" binding:"required" desc:"密码"`
	}
	swaggerGenerator.AddAPI(gb.SwaggerAPIInfo{
		Path:        "/user",
		Method:      "post", //大小写不敏感
		Summary:     "创建用户",
		Description: "创建用户",
		Tags:        []string{"user"},
		Request:     CreateUser{},
	})
	swaggerGenerator.Generate()
}
