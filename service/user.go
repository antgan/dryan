package service

import (
	"context"
	"dryan/dao"
	"dryan/model/do"
	"dryan/util"
)

func QueryUserById(ctx context.Context, userId string) (*do.User, error) {
	objId, _ := util.StringToObjectId(userId)
	var user *do.User
	err := dao.UserOp.FindById(ctx, objId, &user)
	return user, err
}
