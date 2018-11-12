package voucher

import (
	oAuth "ehelp/o/user/auth"
	oVou "ehelp/o/voucher"
	"ehelp/x/rest"
	"ehelp/x/web"
	//"errors"
	"github.com/gin-gonic/gin"
)

type VoucherServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewVoucherServer(parent *gin.RouterGroup, name string) {
	var s = VoucherServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("/list", s.handleList)
	//s.POST("/create", s.handleCreate)
	s.GET("/get", s.handleGet)
}

// func (s *VoucherServer) handleCreate(ctx *gin.Context) {
// 	var body *oVou.Voucher
// 	ctx.BindJSON(&body)
// 	var vou, err = body.Create()
// 	rest.AssertNil(rest.BadRequestValid(err))
// 	oVou.VoucherCache = append(oVou.VoucherCache, vou)
// 	s.SendData(ctx, vou)
// }

func (s *VoucherServer) handleList(ctx *gin.Context) {
	oAuth.GetCusFromToken(ctx.Request)
	s.SendData(ctx, oVou.VoucherCache)
}

func (s *VoucherServer) handleGet(ctx *gin.Context) {
	oAuth.GetCusFromToken(ctx.Request)
	var vouId = ctx.Request.URL.Query().Get("code")
	var typeWork = web.MustGetInt64("service_type", ctx.Request.URL.Query())
	var vou, err = oVou.GetVoucherByID(vouId)
	rest.AssertNil(err)
	vou, err = vou.Validate(int(typeWork))
	rest.AssertNil(err)
	s.SendData(ctx, vou)
}
