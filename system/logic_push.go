package system

import (
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/push_token"
	"ehelp/setting"
	"ehelp/x/fcm"
)

func GetSendPushWork() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var res = int(ord.HourStartItem) - common.GetTimeNowVietNam().Hour()
		if !ord.IsUsedBf && setting.SettingSys.AboutHourStartWork == res && (ord.Status == common.ORDER_STATUS_ACCEPTED || ord.Status == common.ORDER_STATUS_WORKING) {
			ords = append(ords, ord)
		}
	}
	return
}

func GetSendPushWorkEnd() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var val = ord.HourEndItem - common.HourMinute()
		if !ord.IsUsedEnd && ord.StatusItem == common.ITEM_ORDER_STATUS_WORKING && val >= 0 && float64(val) <= setting.SettingSys.AboutMinuteFinishWork {
			ords = append(ords, ord)
		}
	}
	return
}

func checkAndSendPushBf() {
	var ordPushWork = GetSendPushWork()
	if len(ordPushWork) > 0 {
		for _, ord := range ordPushWork {
			var body = "Thời gian: "
			var notify = fcm.FmcMessage{
				Title: "Hôm nay có lịch làm việc!",
				Body: body + common.ConvertF32ToString(ord.HourStartItem) +
					".\nĐịa chỉ: " + ord.AddressLoc.Address +
					".\nVui lòng đến đúng giờ!",
			}
			var pushs, _ = push_token.GetPushsUserId(ord.EmpID)
			sendNotify(notify, []string{ord.EmpID}, "", true, pushs, ord.ID, ord.Status)
			ord.IsUsedBf = true
		}

	}
}

func checkAndSendPushEnd() {
	var ordPushWork = GetSendPushWorkEnd()
	if len(ordPushWork) > 0 {
		for _, ord := range ordPushWork {
			var body = "Thời gian: "
			var notify = fcm.FmcMessage{
				Title: "Sắp hết giờ!",
				Body:  body + common.ConvertF32ToString(ord.HourEndItem) + " sẽ đủ số giờ làm. Hãy bấm kết thúc!",
			}
			var pushs, _ = push_token.GetPushsUserId(ord.EmpID)
			sendNotify(notify, []string{ord.EmpID}, "", true, pushs, ord.ID, ord.Status)
			ord.IsUsedEnd = true
		}
	}
}

func getAutoMissed() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var val = ord.HourStartItem - common.HourMinute()
		if !ord.IsUsedMissed && ord.Status == common.ORDER_STATUS_BIDDING && (val >= 0 && val <= float32(setting.SettingSys.TimeHourHiddenOrder)) {
			ords = append(ords, ord)
		}
	}
	return
}

func SendPushAndChangeMissed() {
	var ordMisseds = getAutoMissed()
	var ordIds = make([]string, len(ordMisseds))
	for i, ord := range ordMisseds {
		ordIds[i] = ord.ID
		var body = "Đơn tại: "
		var notify = fcm.FmcMessage{
			Title: "Đơn được hủy!",
			Body: body + ord.AddressLoc.Address + ".\nThời gian " +
				common.ConvertF32ToString(ord.HourStartItem) +
				" sẽ được hủy do không tìm được người làm.\nQuý khách vui lòng lên lại đơn để tìm người giúp việc!",
		}
		var pushs, _ = push_token.GetPushsUserId(ord.CusID)
		sendNotify(notify, nil, ord.CusID, false, pushs, ord.ID, common.ORDER_STATUS_OPEN)
		ord.IsUsedMissed = true
		//delete(CacheOrderByDay.Orders, ord.ID)
	}
	order.UpdateStatusByIds(ordIds, common.ORDER_STATUS_OPEN)
}
