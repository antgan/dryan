package service

import (
	"context"
	"dryan/dao"
	"dryan/model/do"
	"dryan/model/vo"
	"dryan/util"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func QueryUserById(ctx context.Context, userId string) (*do.User, error) {
	objId, _ := util.StringToObjectId(userId)
	var user *do.User
	err := dao.UserOp.FindById(ctx, objId, &user)
	return user, err
}

func QueryUserByName(ctx context.Context, name string) (*vo.UserVO, error) {
	q := bson.M{"name": name}
	var user *do.User
	err := dao.UserOp.FindOne(ctx, &user, q, nil, nil, 0)
	if err != nil {
		logutil.Errorf("query user by name failed, name:%s, err:%v", name, err)
		return nil, err
	}

	userVO := &vo.UserVO{
		Id:         user.Id.Hex(),
		Name:       user.Name,
		Type:       user.Type,
		CreateTime: user.CreateTime,
	}

	return userVO, nil
}
