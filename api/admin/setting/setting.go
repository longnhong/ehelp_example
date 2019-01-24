package setting

import (
	"ehelp/o/setting"
	"ehelp/x/rest"
	"g/x/web"
	"github.com/gin-gonic/gin"
)

type SettingServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewSettingServer(parent *gin.RouterGroup, name string) {
	var s = &SettingServer{
		RouterGroup: parent.Group(name),
	}

	s.POST("/update", s.handleUpdate)
	s.POST("/create", s.handleCreate)
	s.GET("/list", s.handleGetList)

}

func (s *SettingServer) handleUpdate(ctx *gin.Context) {
	var t *setting.Setting
	web.AssertNil(ctx.ShouldBindJSON(&t))
	var tk, err = t.Update()
	web.AssertNil(err)
	s.SendData(ctx, tk)
}

func (s *SettingServer) handleCreate(ctx *gin.Context) {
	var t *setting.Setting
	web.AssertNil(ctx.ShouldBindJSON(&t))
	var tk, err = t.Create()
	web.AssertNil(err)
	s.SendData(ctx, tk)
}

func (s *SettingServer) handleGetList(ctx *gin.Context) {
	var lst, err = setting.GetAll()
	web.AssertNil(err)
	s.SendData(ctx, lst)
}
