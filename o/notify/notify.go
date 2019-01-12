package notify

import (
	"ehelp/common"
	"ehelp/x/db/mongodb"
)

type Notify struct {
	mongodb.BaseModel `bson:",inline"`
	EmpIDs            []string           `bson:"emps" json:"emps"`
	CusID             string             `bson:"cus_id" json:"cus_id"`
	OrderID           string             `bson:"order_id" json:"order_id"`
	Body              string             `bson:"body" json:"body"`
	Title             string             `bson:"title" json:"title"`
	StatusOder        common.OrderStatus `bson:"status_oder" json:"status_oder"`
}

var notifyTable = mongodb.NewTable("notify", "nf", 25)

func (noti *Notify) Create() (*Notify, error) {
	err := notifyTable.Create(noti)
	return noti, err
}
