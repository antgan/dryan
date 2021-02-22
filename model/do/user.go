package do

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"name"`
	Type       string        `bson:"type"` //official官方；director董事
	CreateTime time.Time     `bson:"create_time"`
}
