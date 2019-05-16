package system

import (
	"ehelp/common"
	"ehelp/o/order"
	"ehelp/o/push_token"
	"ehelp/o/user/customer"
	"ehelp/o/user/employee"
	"ehelp/setting"
	"ehelp/x/fcm"
)

func GetSendPushWork() (ords []*order.OrderCheckChange) {
	ords = make([]*order.OrderCheckChange, 0)
	for _, ord := range CacheOrderByDay.Orders {
		var res = int(ord.HourStartItem) - common.GetTimeNowVietNam().Hour()
		if !ord.IsUsedBf && setting.SettingSys.AboutHourGoWork == res && (ord.Status == common.ORDER_STATUS_ACCEPTED || ord.Status == common.ORDER_STATUS_WORKING) {
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
			var body string
			var notify = fcm.FmcMessage{}
			var lang common.Lang
			var user, err = employee.GetByID(ord.EmpID)
			if err != nil {
				lang = common.LANG_VI
			} else {
				lang = user.Lang
				if len(lang) == 0 {
					lang = common.LANG_VI
				}
			}
			switch lang {
			case common.LANG_VI:
				body = TimeVI
				notify = fcm.FmcMessage{
					Title: TitlePushWorkVI,
					Body: body + common.ConvertF32ToString(ord.HourStartItem) +
						Body1PushWorkVi + ord.AddressLoc.Address + Body2PushWorkVi,
				}
			case common.LANG_EN:
				body = TimeEN
				notify = fcm.FmcMessage{
					Title: TitlePushWorkEN,
					Body: body + common.ConvertF32ToString(ord.HourStartItem) +
						Body1PushWorkEN + ord.AddressLoc.Address + Body2PushWorkEN,
				}
			case common.LANG_CN:
				body = TimeCN
				notify = fcm.FmcMessage{
					Title: TitlePushWorkCN,
					Body: body + common.ConvertF32ToString(ord.HourStartItem) +
						Body1PushWorkCN + ord.AddressLoc.Address + Body2PushWorkCN,
				}
			}

			// var body = "Thời gian: "
			// var notify = fcm.FmcMessage{
			// 	Title: "Hôm nay có lịch làm việc!",
			// 	Body: body + common.ConvertF32ToString(ord.HourStartItem) +
			// 		".\nĐịa chỉ: " + ord.AddressLoc.Address +
			// 		".\nVui lòng đến đúng giờ!",
			// }
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
			var body string
			var notify = fcm.FmcMessage{}
			var lang common.Lang
			var user, err = employee.GetByID(ord.EmpID)
			if err != nil {
				lang = common.LANG_VI
			} else {
				lang = user.Lang
				if len(lang) == 0 {
					lang = common.LANG_VI
				}
			}
			switch lang {
			case common.LANG_VI:
				body = TimeVI
				notify = fcm.FmcMessage{
					Title: TitlePushNearFinishVI,
					Body:  body + common.ConvertF32ToString(ord.HourEndItem) + BodyPushNearFinishVI,
				}
			case common.LANG_EN:
				body = TimeEN
				notify = fcm.FmcMessage{
					Title: TitlePushWorkEN,
					Body:  body + common.ConvertF32ToString(ord.HourEndItem) + BodyPushNearFinishEN,
				}
			case common.LANG_CN:
				body = TimeCN
				notify = fcm.FmcMessage{
					Title: TitlePushWorkCN,
					Body:  body + common.ConvertF32ToString(ord.HourEndItem) + BodyPushNearFinishCN,
				}
			}
			// var body = "Thời gian: "
			// var notify = fcm.FmcMessage{
			// 	Title: "Sắp hết giờ!",
			// 	Body:  body + common.ConvertF32ToString(ord.HourEndItem) + " sẽ đủ số giờ làm. Hãy bấm kết thúc!",
			// }
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
		var notify = fcm.FmcMessage{}
		var lang common.Lang
		var user, err = customer.GetByID(ord.CusID)
		if err != nil {
			lang = common.LANG_VI
		} else {
			lang = user.Lang
			if len(lang) == 0 {
				lang = common.LANG_VI
			}
		}
		switch lang {
		case common.LANG_VI:
			notify = fcm.FmcMessage{
				Title: TitleChangeMissedVI,
				Body: Body1ChangeMissedVI + ord.AddressLoc.Address + TimeVI +
					common.ConvertF32ToString(ord.HourStartItem) + Body2ChangeMissedVI,
			}
		case common.LANG_EN:
			notify = fcm.FmcMessage{
				Title: TitleChangeMissedEN,
				Body: Body1ChangeMissedEN + ord.AddressLoc.Address + TimeEN +
					common.ConvertF32ToString(ord.HourStartItem) + Body2ChangeMissedEN,
			}
		case common.LANG_CN:
			notify = fcm.FmcMessage{
				Title: TitleChangeMissedCN,
				Body: Body1ChangeMissedCN + ord.AddressLoc.Address + TimeCN +
					common.ConvertF32ToString(ord.HourStartItem) + Body2ChangeMissedCN,
			}
		}

		// var body = "Đơn tại: "
		// var notify = fcm.FmcMessage{
		// 	Title: "Đơn được hủy!",
		// 	Body: body + ord.AddressLoc.Address + ".\nThời gian " +
		// 		common.ConvertF32ToString(ord.HourStartItem) +
		// 		" sẽ được hủy do không tìm được người làm.\nQuý khách vui lòng lên lại đơn để tìm người giúp việc!",
		// }
		var pushs, _ = push_token.GetPushsUserId(ord.CusID)
		sendNotify(notify, nil, ord.CusID, false, pushs, ord.ID, common.ORDER_STATUS_OPEN)
		ord.IsUsedMissed = true
		//delete(CacheOrderByDay.Orders, ord.ID)
	}
	order.UpdateStatusByIds(ordIds, common.ORDER_STATUS_OPEN)
}
