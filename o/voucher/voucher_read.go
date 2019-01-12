package voucher

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

func GetListVoucher() (lstVoucher []*Voucher, err error) {
	var timeNow = time.Now().Unix()
	var match = bson.M{
		"active": true,
		"start_time": bson.M{
			"$lte": timeNow,
		},
		"end_time": bson.M{
			"$gte": timeNow,
		},
	}
	return lstVoucher, VoucherTable.FindWhere(match, &lstVoucher)
}

func GetVoucherByCodeTime(code string) (vou *Voucher, err error) {
	var timeNow = time.Now().Unix()
	var match = bson.M{
		"code":   code,
		"active": true,
		"start_time": bson.M{
			"$lte": timeNow,
		},
		"end_time": bson.M{
			"$gte": timeNow,
		},
	}
	return vou, VoucherTable.FindOne(match, &vou)
}
func GetVoucherCode(code string) (vou *Voucher, err error) {
	var match = bson.M{
		"code":   code,
		"active": true,
	}
	return vou, VoucherTable.FindOne(match, &vou)
}
