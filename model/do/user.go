package do

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"name"`
	Account    string        `bson:"account"`
	Password   string        `bson:"password"`
	CreateTime time.Time     `bson:"create_time"`
}
