package db

import (
	"dryan/common"
	"errors"
	logutil "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

type MgoSession struct {
	*mgo.Session
	dbname   string
	user     string
	pass     string
	sourcedb string
}

type MgoDB struct {
	*mgo.Database
}

var dbConns = make(map[string]interface{}, 0)

func Mongo(key string) (mongo *MgoSession, err error) {
	if v, ok := dbConns[key]; ok {
		return v.(*MgoSession), nil
	} else {
		return nil, errors.New("connection not exists")
	}
}

func init() {
	dbs := common.Config.DATABASES
	if dbs == nil {
		logutil.Info("No database config")
		return
	}

	for _, _db := range dbs {
		if _db.TYPE == "mongodb" {
			m, err := initMongoDB(&_db)
			if err != nil {
				logutil.Info("Init mongodb failed.")
			} else {
				dbConns[_db.KEY] = m
				logutil.Info("Init mongodb success.")
			}
		}
	}
}
