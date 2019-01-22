package push_token

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*PushToken, error) {
	var auth *PushToken
	return auth, PushTokenTable.FindOne(bson.M{
		"_id":       id,
		"is_revoke": false,
	}, &auth)
}

func GetPushsUserId(userId string) ([]string, error) {
	var auth *PushToken
	err := PushTokenTable.FindOne(bson.M{
		"user_id":   userId,
		"is_revoke": false,
	}, &auth)
	var p string
	if auth != nil {
		p = auth.PushToken
	}
	return []string{p}, err
}

func GetPushsUserIds(userIds []string) ([]string, error) {
	var pushs []string
	var err = PushTokenTable.Find(bson.M{
		"user_id":   bson.M{"$in": userIds},
		"is_revoke": false,
	}).Distinct("push_token", &pushs)
	fmt.Println(len(pushs))
	return pushs, err
}
