package examples

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/loveyu233/gb"
)

func TestParamA(t *testing.T) {
	gb.PublicRoutes = append(gb.PublicRoutes, func(group *gin.RouterGroup) {
		group.GET("/a/:id", func(c *gin.Context) {
			type Req struct {
				Name     string `json:"name" binding:"required" label:"姓名"`
				Email    string `json:"email" binding:"required,email" label:"邮箱"`
				Phone    string `json:"phone" binding:"required,phone" label:"手机号"`
				Age      int    `json:"age" binding:"required,gte=0,lte=120" label:"年龄"`
				Password string `json:"password" binding:"required,min=6" label:"密码"`
				//ID int64 `uri:"id" binding:"required"`
			}
			//var (
			//	req = new(Req)
			//)
			//
			//if err := c.BindJSON(req); err != nil {
			//	gb.ResponseParamError(c, err)
			//	return
			//}
			//c.Param("")
			//c.Query("key")
			//ageInt, err := gb.QueryRequired[int](c, "age")
			ageInt, err := gb.GetGinPathRequired[int](c, "id")
			if err != nil {
				gb.ResponseParamError(c, err)
			}
			//optional, err := gb.QueryOptional[int](c, "phone")
			//if err != nil {
			//	gb.ResponseParamError(c, err)
			//}
			//size, err := gb.QueryOptionalWithDefault(c, "size", 10)
			//if err != nil {
			//	gb.ResponseParamError(c, err)
			//}
			//isActive, err := gb.QueryOptional[bool](c, "is_active")
			//if err != nil {
			//	gb.ResponseParamError(c, err)
			//}
			gb.ResponseSuccess(c, map[string]any{
				"ageInt": ageInt,
				//"optional": optional,
				//"size":     size,
				//"isActive": isActive,
			})
		})
	})
	gb.InitHTTPServerAndStart(gb.DefaultGinTokenConfig, "127.0.0.1:9999")
}
