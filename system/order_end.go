package system

import (
	"ehelp/cache"
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/order_hst"
	"ehelp/o/push_token"
	oAuth "ehelp/o/user/auth"
	"ehelp/x/fcm"
	"fmt"
)

func canceled(ord *order.Order) {
	CreateOrderHst(ord.CusID, ord.ID, ord.ServiceWorks, common.ORDER_STATUS_CANCELED)
}
func accepted(ord *order.Order) {
	fmt.Printf("========= EmpID: "+ord.EmpID+" ID: "+ord.ID, ord.ServiceWorks)
	CreateOrderHst(ord.EmpID, ord.ID, ord.ServiceWorks, common.ORDER_STATUS_ACCEPTED)
	var emp, _ = cache.GetEmpID(ord.EmpID)
	fmt.Println("==== KHÁCH" + ord.CusID)
	pushs, _ := push_token.GetPushsUserId(ord.CusID) // phải distinic
	if len(pushs) > 0 {
		// var noti = fcm.FmcMessage{
		// 	Title: "Đã có người nhận!",
		// 	Body:  emp.FullName + " đã nhận đơn"}
		var noti = fcm.FmcMessage{}
		if len(emp.Lang) == 0 {
			emp.Lang = common.LANG_VI
		}
		switch emp.Lang {
		case common.LANG_VI:
			noti.Title = TitleAccepVI
			noti.Body = BodyAccepVI
		case common.LANG_EN:
			noti.Title = TitleAccepEN
			noti.Body = BodyAccepEN
		case common.LANG_CN:
			noti.Title = TitleAccepCN
			noti.Body = BodyAccepCN
		}
		sendNotify(noti, nil, ord.CusID, false, pushs, ord.ID, common.ORDER_STATUS_ACCEPTED)
	}
}
func bidding(ord *order.Order) {
	var langs = []common.Lang{common.LANG_VI, common.LANG_EN, common.LANG_CN}
	for _, lang := range langs {
		var empIds, _ = oAuth.GetListEmpVsOrderBidding(ord.ServiceWorks, ord.AddressLoc.Address, string(lang))
		var pushs, err2 = push_token.GetPushsUserIds(empIds) // phải distinic
		logAction.Errorf("push_token.GetPushsUserId", err2)
		// var noti = fcm.FmcMessage{
		// 	Title: "Có việc mới!",
		// 	Body:  "Công việc tại " + ord.AddressLoc.Address,
		// }
		var noti = fcm.FmcMessage{}
		switch lang {
		case common.LANG_VI:
			noti.Title = TitleAccepVI
			noti.Body = BodyAccepVI
		case common.LANG_EN:
			noti.Title = TitleAccepEN
			noti.Body = BodyAccepEN
		case common.LANG_CN:
			noti.Title = TitleAccepCN
			noti.Body = BodyAccepCN
		}
		sendNotify(noti, empIds, "", true, pushs, ord.ID, common.ORDER_STATUS_BIDDING)
		CreateOrderHst(ord.CusID, ord.ID, ord.ServiceWorks, common.ORDER_STATUS_BIDDING)
	}

}

func working(ord *order.Order, itemOrder *common.DayWeek) {
	CreateItemOrderHst(itemOrder.HourDay, float32(ord.PriceEnd),
		ord.CusID, ord.EmpID, ord.ServiceWorks, ord.ID, common.ORDER_STATUS_WORKING, itemOrder.IdItem,
		common.ITEM_ORDER_STATUS_WORKING, itemOrder.HourStart, itemOrder.HourEnd)
	var pushs, _ = push_token.GetPushsUserId(ord.CusID) // phải distinic
	var emp, _ = cache.GetEmpID(ord.EmpID)
	// var noti = fcm.FmcMessage{
	// 	Title: "Bắt đầu làm việc!",
	// 	Body:  empOrd.FullName + " vừa bắt đầu làm việc!"}
	var noti = fcm.FmcMessage{}
	if len(emp.Lang) == 0 {
		emp.Lang = common.LANG_VI
	}
	switch emp.Lang {
	case common.LANG_VI:
		noti.Title = TitleWorkingVI
		noti.Body = BodyWorkingVI
	case common.LANG_EN:
		noti.Title = TitleWorkingEN
		noti.Body = BodyWorkingEN
	case common.LANG_CN:
		noti.Title = TitleWorkingCN
		noti.Body = BodyWorkingCN
	}
	sendNotify(noti, nil, ord.CusID, false, pushs, ord.ID, common.ORDER_STATUS_WORKING)
}

func finished(ord *order.Order, itemOrder *common.DayWeek) {
	CreateItemOrderHst(itemOrder.HourDay, float32(ord.PriceEnd),
		ord.CusID, ord.EmpID, ord.ServiceWorks, ord.ID, common.ORDER_STATUS_WORKING, itemOrder.IdItem,
		common.ITEM_ORDER_STATUS_FINISHED, itemOrder.HourStart, itemOrder.HourEnd)
	var countCusNew, _ = order.CheckCustomerNewOfEmployee(ord.CusID, ord.EmpID)
	oAuth.UpdateCusNewAndHour(ord.EmpID, countCusNew, itemOrder.HourDay)
	var pushs, _ = push_token.GetPushsUserId(ord.CusID) // phải distinic
	var emp, _ = cache.GetEmpID(ord.EmpID)
	if ord.Status == common.ORDER_STATUS_FINISHED {
		// var noti = fcm.FmcMessage{
		// 	Title: "Công việc đã hoàn thành!",
		// 	Body:  empOrd.FullName + " đã hoàn thành đầy đủ đơn mà bạn đặt! Lên đơn mới nếu muốn tìm người giúp việc!"}

		var noti = fcm.FmcMessage{}
		if len(emp.Lang) == 0 {
			emp.Lang = common.LANG_VI
		}
		switch emp.Lang {
		case common.LANG_VI:
			noti.Title = TitleFinishVI
			noti.Body = BodyFinishAllVI
		case common.LANG_EN:
			noti.Title = TitleFinishEN
			noti.Body = BodyFinishAllEN
		case common.LANG_CN:
			noti.Title = TitleFinishCN
			noti.Body = BodyFinishAllCN
		}
		sendNotify(noti, nil, ord.CusID, false, pushs, ord.ID, common.ORDER_STATUS_FINISHED)
	} else {
		// var noti = fcm.FmcMessage{
		// 	Title: "Công việc đã hoàn thành!",
		// 	Body:  empOrd.FullName + " đã hoàn thành việc ngày hôm nay!"}
		var noti = fcm.FmcMessage{}
		if len(emp.Lang) == 0 {
			emp.Lang = common.LANG_VI
		}
		switch emp.Lang {
		case common.LANG_VI:
			noti.Title = TitleFinishVI
			noti.Body = BodyFinishItemVI
		case common.LANG_EN:
			noti.Title = TitleFinishEN
			noti.Body = BodyFinishItemEN
		case common.LANG_CN:
			noti.Title = TitleFinishCN
			noti.Body = BodyFinishItemCN
		}
		sendNotify(noti, nil, ord.CusID, false, pushs, ord.ID, common.ORDER_STATUS_FINISHED)
	}
}

func CreateOrderHst(userId string, orderID string, services []string, statusOrder common.OrderStatus) {
	var ordHst = order_hst.OrderHST{
		CusId:       userId,
		Services:    services,
		OrderId:     orderID,
		OrderStatus: statusOrder,
		ItemStatus:  common.ITEM_ORDER_STATUS_NEW,
	}
	ordHst.CrateOrderHistory()
}
func CreateItemOrderHst(itemHour float32, itemMoney float32, cusId string, empId string, services []string, orderID string, statusOrder common.OrderStatus, itemId string, statusItem common.ItemOrderStatus, startDay float32, endDay float32) {
	var ordHst = order_hst.OrderHST{
		CusId:        cusId,
		EmpId:        empId,
		Services:     services,
		ItemId:       itemId,
		ItemStatus:   statusItem,
		OrderId:      orderID,
		OrderStatus:  statusOrder,
		HourDay:      itemHour,
		MoneyDay:     itemMoney,
		StartWorkDay: startDay,
		EndWorkDay:   endDay,
	}
	ordHst.CrateOrderHistory()
}
