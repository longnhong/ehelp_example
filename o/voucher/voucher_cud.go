package voucher

import (
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
)

func (t *Voucher) Create() (*Voucher, error) {
	if err := validate.Struct(t); err != nil {
		glog.Error(err)
		return nil, err
	}
	return t, VoucherTable.CreateUnique(bson.M{"code": t.Code}, t)
}

func (t *Voucher) Update() error {
	return VoucherTable.UpdateByID(t.ID, t)
}

func UpdateCount(code []string) error {
	var match = bson.M{
		"code": bson.M{"$in": code},
	}
	var _, err = VoucherTable.UpdateAll(match, bson.M{"$inc": bson.M{"count": 1}})
	return err
}
