package auth

import (
	"ehelp/x/db/mongodb"
	"gopkg.in/mgo.v2/bson"
)

type Auth struct {
	mongodb.BaseModel `bson:",inline"`
	Role              string `bson:"role" json:"role"`
	UserID            string `bson:"user_id" json:"user_id"`
	//	IsRevoke          bool   `bson:"is_revoke" json:"is_revoke"`
}

var AuthTable = mongodb.NewTable("auth", "k", 80)

func Create(userID string, role string) (*Auth, error) {
	var a = &Auth{}
	a.UserID = userID
	a.Role = role
	err := AuthTable.Create(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func UpdateRevoke(id string, revoke bool) error {
	err := AuthTable.UpdateSetByID(id, bson.M{
		"is_revoke": true,
	})
	return err
}

func GetByID(id string) (*Auth, error) {
	var auth *Auth
	return auth, AuthTable.FindByID(id, &auth)
}
