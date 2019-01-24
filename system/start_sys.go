package system

import (
	//"ehelp/cache"

	"ehelp/o/order"
	"ehelp/o/setting"
	"ehelp/o/voucher"

	"ehelp/x/mlog"
	"ehelp/x/rest"
	"fmt"
	"time"
)

var logSys = mlog.NewTagLog("System")

/*Update setting*/
func updateSetting() {
	setting.UpdateSetting()
}

/*Get Order ID*/
func GetOrderID(id string) (*order.Order, bool) {
	if val, ok := CacheOrderByDay.Orders[id]; ok {
		return val.Order, true
	}
	var ord, err = order.GetOrderById(id)
	rest.AssertNil(err)
	// if ord != nil {
	// 	var weekDayNow = common.GetTimeNowVietNam().Weekday()
	// 	for _, itemWork := range ord.DayWeeks {
	// 		if common.ConvertTimeEpochToWeek(itemWork.DateIn) == weekDayNow {
	// 			var or = order.OrderCheckChange{
	// 				Order:         ord,
	// 				HourStartItem: itemWork.HourStart,
	// 				HourEndItem:   itemWork.HourEnd,
	// 				StatusItem:    itemWork.Status,
	// 			}
	// 			CacheOrderByDay.Orders[id] = &or
	// 			break
	// 		}
	// 	}
	// }
	return ord, false
}

func Launch() {
	//cache.SetCacheCus()
	//cache.SetCacheEmp()
	SetCacheOrderDay()
	//voucher.VoucherCache, _ = voucher.GetListVoucher()
	updateSetting()
	fmt.Println("SỐ VOUS", voucher.VoucherCache)
	go startCache(CacheOrderByDay)
	go everyCheckSend(CacheOrderByDay)
}

func startCache(c *CacheOrderWorker) {
	fmt.Println("Vào startCache")
	for {
		select {
		case action := <-c.OrderUpdate:
			c.OrderWorking(action)
		}
	}
}

func everyCheckSend(c *CacheOrderWorker) {
	everyMinute := time.Tick(5 * time.Minute)
	for {
		select {
		case <-everyMinute:
			fmt.Print("everyMinute")
			checkAndSendPushEnd()
			checkAndSendPushBf()
			SendPushAndChangeMissed()
		}
	}
}
