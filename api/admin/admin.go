package admin

import (
	"ehelp/api/admin/report"
	"ehelp/api/admin/service"
	"ehelp/api/admin/setting"
	"ehelp/api/admin/user"
	"ehelp/middleware"
	usr "ehelp/o/admin/user"
	"ehelp/o/auth"
	"ehelp/o/order_hst"
	oVou "ehelp/o/voucher"
	"ehelp/x/rest"
	"github.com/gin-gonic/gin"
)

type AdminServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

//comment log
func NewAdminServer(parent *gin.RouterGroup, name string) *AdminServer {
	var s = AdminServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/refresh_config", s.handleRefresh)
	s.POST("/signin", s.handleSignin)
	s.GET("/super_admin", s.handleExistSuperAdmin)
	s.POST("/register", s.handleRegister)
	s.GET("/order-history", s.handleOrderHistory)
	s.Use(middleware.MustBeSuperAdmin())
	service.NewServiceServer(s.RouterGroup, "service")
	report.NewReportServer(s.RouterGroup, "report")
	user.NewUserServer(s.RouterGroup, "user")
	setting.NewSettingServer(s.RouterGroup, "setting")
	return &s
}

func (s *AdminServer) handleRefresh(ctx *gin.Context) {
	var body = struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{}
	ctx.BindJSON(&body)
	_, err := usr.GetByUNamePwd(body.Username, body.Password, body.Role)
	rest.AssertNil(err)
	res, err := oVou.GetListVoucher()
	rest.AssertNil(err)
	oVou.VoucherCache = res
	s.SendData(ctx, oVou.VoucherCache)
}

func (s *AdminServer) handleOrderHistory(ctx *gin.Context) {
	var res, err = order_hst.GetOrderHistory()
	rest.AssertNil(err)
	s.SendData(ctx, res)
}
func (s *AdminServer) handleExistSuperAdmin(ctx *gin.Context) {
	s.SendData(ctx, usr.GetSuperUser())
}
func (s *AdminServer) handleRegister(ctx *gin.Context) {
	var u *usr.Admin
	ctx.BindJSON(&u)
	rest.AssertNil(u.Create())
	auth, err := auth.Create(u.ID, string(u.Role))
	rest.AssertNil(err)
	s.SendData(ctx, map[string]interface{}{
		"access_token": auth.ID,
	})
}
func (s *AdminServer) handleSignin(ctx *gin.Context) {
	var loginInfo = struct {
		UName    string `json:"uname"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{}
	ctx.BindJSON(&loginInfo)
	u, err := usr.GetByUNamePwd(loginInfo.UName, loginInfo.Password, loginInfo.Role)
	rest.AssertNil(err)
	auth, err := auth.Create(u.ID, string(u.Role))
	rest.AssertNil(err)
	s.SendData(ctx, map[string]interface{}{
		"access_token": auth.ID,
		"role":         auth.Role,
	})
}
