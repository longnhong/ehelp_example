package voucher

import (
	"ehelp/x/db/mongodb"
	validator "gopkg.in/go-playground/validator.v9"
)

type Voucher struct {
	mongodb.BaseModel `bson:",inline"`
	Title             string  `bson:"title" json:"title" validate:"required"`
	Description       string  `bson:"description" json:"description" validate:"required"`
	Code              string  `bson:"code" json:"code" validate:"required"`
	AutoActive        bool    `bson:"auto_active" json:"auto_active" `
	ServiceType       int     `bson:"service_type" json:"service_type" validate:"required"`
	Value             float32 `bson:"value" json:"value"`
	ValueRatio        float32 `bson:"value_ratio" json:"value_ratio"`
	Active            bool    `bson:"active" json:"active" validate:"required"`
	StartTime         int64   `bson:"start_time" json:"start_time" validate:"required"`
	EndTime           int64   `bson:"end_time" json:"end_time" validate:"required"`
	Quantity          int     `bson:"quantity" json:"quantity"`
	Count             int     `bson:"count" json:"count"`
}

var VoucherTable = mongodb.NewTable("voucher", "vou", 20)

var VoucherCache = make([]*Voucher, 0)

var validate = validator.New()
