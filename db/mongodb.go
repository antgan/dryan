package db

import (
	"dryan/common"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"sync"
)

var (
	lock = sync.Mutex{}
)

func initMongoDB(dbconf *common.DataBase) (session *MgoSession, err error) {
	lock.Lock()
	defer lock.Unlock()
	logutil.Infof("Ready to connect to mongodb. key: %v", dbconf.KEY)
	dbSession, err := mgo.Dial(dbconf.HOST)
	if err != nil {
		logutil.Infof("Connect mongodb failed: key: %v, err: %v ", dbconf.KEY, err)
		return nil, err
	}
	dbSession.SetPoolLimit(dbconf.Ext("maxPoolSize", 256).(int))
	// 2: Strong模式，牺牲速度保证一致性
	dbSession.SetMode(mgo.Mode(dbconf.Ext("mode", 2).(int)), true)

	if dbconf.USER != "" {
		dbSession.Login(&mgo.Credential{
			Username: dbconf.USER,
			Password: dbconf.PASSWORD,
			Source:   dbconf.Ext("authSource", dbconf.NAME).(string),
		})
	}

	dbSession.Refresh()
	logutil.Infof("Connect mongodb success. key: %s", dbconf.KEY)
	return &MgoSession{
		Session:  dbSession,
		dbname:   dbconf.NAME,
		user:     dbconf.USER,
		pass:     dbconf.PASSWORD,
		sourcedb: dbconf.Ext("authSource", dbconf.NAME).(string),
	}, nil
}

func (s *MgoSession) DB() (database *MgoDB, err error) {
	_s := s.Copy()
	db := _s.DB(s.dbname)
	return &MgoDB{db}, nil
}

func (d *MgoDB) Close() {
	d.Session.Close()
}
