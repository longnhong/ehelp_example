package fcm

import (
	fcm "github.com/NaySoftware/go-fcm"
	"time"
)

const (
	RESPONSE_FAIL = "fail"
)

type FcmClient struct {
	*fcm.FcmClient
}

type FmcMessage struct {
	Title string      `json:"title,omitempty"`
	Body  string      `json:"body,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func NewFCM(serverKey string) *FcmClient {
	return &FcmClient{
		FcmClient: fcm.NewFcmClient(serverKey),
	}
}

type dataMsg struct {
	Content interface{} `json:"content,omitempty"`
	Time    int64       `json:"time,omitempty"`
	Unseen  int         `json:"unseen,omitempty"`
}

func (f *FcmClient) SendToMany(ids []string, data FmcMessage) (error, string) {
	var noti = fcm.NotificationPayload{
		Title: data.Title,
		Body:  data.Body,
		Sound: "ting.wav",
	}
	var dataMs = dataMsg{
		Content: data.Data,
		Time:    time.Now().Unix(),
	}
	f.NewFcmRegIdsMsg(ids, map[string]interface{}{
		"data": dataMs})
	f.SetNotificationPayload(&noti)
	f.SetContentAvailable(true)
	// f.SetMsgData(map[string]interface{}{"notify": data.Data})
	status, err := f.Send()
	if err != nil {
		return err, RESPONSE_FAIL
	}
	return nil, status.Err
}

func (f *FcmClient) SendToOne(id string, data FmcMessage) (error, string) {
	return f.SendToMany([]string{id}, data)
}
