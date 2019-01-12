package system

import (
	"ehelp/common"
	"ehelp/o/notify"
	"ehelp/x/fcm"
)

func sendNotify(fData fcm.FmcMessage, empIDs []string, cusID string, isSendEmp bool, pushs []string, status common.OrderStatus) {
	var noti = notify.Notify{
		Title:      fData.Title,
		Body:       fData.Body,
		EmpIDs:     empIDs,
		CusID:      cusID,
		StatusOder: status,
	}
	noti.Create()
	fData.Data = noti
	if isSendEmp {
		fcm.FcmEmployee.SendToMany(pushs, fData)
	} else {
		fcm.FcmCustomer.SendToMany(pushs, fData)
	}
}
