package api

import (
	"ehelp/api/admin"
	"ehelp/api/auth"
	"ehelp/api/public"
	"ehelp/api/service"
	//"ehelp/api/test_table"
	order "ehelp/api/order_server"
	"ehelp/api/user"
	"ehelp/api/voucher"
	"github.com/gin-gonic/gin"
)

func InitApi(root *gin.RouterGroup) {
	service.NewServiceServer(root, "service")
	user.NewUserServer(root, "user")
	auth.NewAuthServer(root, "auth")
	admin.NewAdminServer(root, "admin")
	order.NewOrderServer(root, "order")
	voucher.NewVoucherServer(root, "voucher")
	auth.NewAuthCustomerServer(root, "customer/auth")
	auth.NewAuthEmployeeServer(root, "employee/auth")
	public.NewPublicServer(root, "public")
	//test_table.NewTableServer(root, "service")
}
