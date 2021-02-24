package service

import (
	"context"
	"dryan/constant"
	"dryan/dao"
	"dryan/model/do"
	"dryan/model/vo"
	"dryan/util"
	"errors"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func QueryUserById(ctx context.Context, userId string) (*do.User, error) {
	objId, _ := util.StringToObjectId(userId)
	var user *do.User
	err := dao.UserOp.FindById(ctx, objId, &user)
	return user, err
}

func QueryUserByName(ctx context.Context, name string) (*vo.UserVO, error) {
	user, err := queryUserByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("not found")
	}

	userVO := &vo.UserVO{
		Id:         user.Id.Hex(),
		Name:       user.Name,
		Type:       user.Type,
		CreateTime: util.FormatTime(user.CreateTime),
	}

	return userVO, nil
}

func queryUserByName(ctx context.Context, name string) (*do.User, error) {
	q := bson.M{"name": name}
	var user *do.User
	err := dao.UserOp.FindOne(ctx, &user, q, nil, nil, 0)
	if err != nil && err != mgo.ErrNotFound {
		logutil.Errorf("query user by name failed, name:%s, err:%v", name, err)
		return nil, err
	}
	return user, nil
}

func AddUser(ctx context.Context, req *vo.AddUserReq) error {
	if req.Type == "" {
		req.Type = constant.DRYAN_USER_TYPE_OFFICIAL
	}
	dbUser, err := queryUserByName(ctx, req.Name)
	if err != nil {
		logutil.Errorf("query user by name failed, name:%v, err:%v", req.Name, err)
		return err
	}
	if dbUser != nil {
		return errors.New("duplicate name")
	}

	err = dao.UserOp.Insert(ctx, &do.User{
		Id:         bson.NewObjectId(),
		Name:       req.Name,
		Type:       req.Type,
		CreateTime: time.Now(),
	})
	if err != nil {
		logutil.Errorf("add user failed, err:%v", err)
		return err
	}
	return nil
}
