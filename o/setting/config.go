package setting

import (
	"ehelp/setting"
	"ehelp/x/db/mongodb"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Setting struct {
	mongodb.BaseModel `bson:",inline"`
	Code              string      `bson:"code" json:"code"`
	Value             interface{} `bson:"value" json:"value"`
	Des               string      `bson:"des" json:"des"`
	Track             []track     `bson:"tracks" json:"tracks"`
}

type track struct {
	Mtime int64       `bson:"mtime" json:"mtime"`
	Value interface{} `bson:"value" json:"value"`
}

var settingTable = mongodb.NewTable("setting", "set", 10)

const (
	TimeHourHiddenOrder   = "time_hour_hidden_order"
	AboutHourStartWork    = "about_hour_start_work"
	AboutMinuteFinishWork = "about_minute_finish_work"
)

func UpdateSetting() {
	var sts, _ = GetAll()
	if sts != nil {
		for _, st := range sts {
			st.updateValue()
		}
	}
}

func (st *Setting) updateValue() {
	if st.Code == TimeHourHiddenOrder {
		setting.SettingSys.TimeHourHiddenOrder = st.Value.(int)
	} else if st.Code == AboutHourStartWork {
		setting.SettingSys.AboutHourStartWork = st.Value.(int)
	} else if st.Code == AboutMinuteFinishWork {
		setting.SettingSys.AboutMinuteFinishWork = st.Value.(float64) / 60
	}
}
func (noti *Setting) Create() (*Setting, error) {
	var tr = track{
		Mtime: time.Now().Unix(),
		Value: noti.Value,
	}
	noti.Track = []track{tr}
	err := settingTable.Create(noti)
	return noti, err
}

func (st *Setting) Update() (*Setting, error) {
	st.Track = append(st.Track, track{
		Mtime: time.Now().Unix(),
		Value: st.Value,
	})
	err := settingTable.UpdateByID(st.ID, st)
	if err == nil {
		st.updateValue()
	}
	return st, err
}

func GetAll() ([]*Setting, error) {
	var sts []*Setting
	return sts, settingTable.FindWhere(bson.M{}, &sts)
}
