package service

import (
	"context"
	"dryan/dao"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func UpdateStockCount(ctx context.Context, userId, itemId string, count int, isAdd bool) error {
	incCount := 0
	if isAdd {
		incCount = count
	} else {
		incCount = count * -1
	}
	q := bson.M{"user_id": userId, "item_id": itemId}
	updateMap := bson.M{"$inc": bson.M{"remain_count": incCount}}
	err := dao.StockOp.Upsert(ctx, q, updateMap)
	if err != nil {
		logutil.Errorf("update stock count failed, err:%v", err)
		return err
	}

	return nil
}
