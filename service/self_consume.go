package service

import (
	"context"
	"dryan/dao"
	"dryan/model/do"
	"dryan/model/vo"
	"dryan/util"
	"errors"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func AddSelfConsume(ctx context.Context, req *vo.AddSelfConsumeReq) error {
	if len(req.Items) <= 0 {
		return errors.New("empty items")
	}
	user, err := QueryUserById(ctx, req.UserId)
	if err != nil || user == nil {
		logutil.Errorf("user not exists, err:%v", err)
		return err
	}

	//插入更新偷吃记录
	for _, item := range req.Items {
		upsertMap := bson.M{
			"$inc": bson.M{"count": item.Count},
			"$set": bson.M{"update_time": time.Now()},
			"$setOnInsert": bson.M{
				"_id":     bson.NewObjectId(),
				"user_id": req.UserId,
				"item_id": item.ItemId,
			},
		}
		q := bson.M{"user_id": req.UserId, "item_id": item.ItemId}
		err = dao.SelfConsumeOp.Upsert(ctx, q, upsertMap)
		if err != nil {
			logutil.Errorf("upsert self consume failed, err:%v", err)
			return err
		}
	}

	//更新库存
	for _, item := range req.Items {
		err = UpdateStockCount(ctx, req.UserId, item.ItemId, item.Count, false)
		if err != nil {
			logutil.Errorf("update stock count failed, err:%v", err)
			return err
		}
	}

	return nil
}

func QueryAllSelfConsume(ctx context.Context, userId string) ([]*vo.SelfConsumeVO, error) {
	q := bson.M{"user_id": userId}
	selfConsumes := make([]*do.SelfConsume, 0)
	err := dao.SelfConsumeOp.Find(ctx, &selfConsumes, q, []string{"-update_time"}, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query self consume failed, err:%v", err)
		return nil, err
	}

	itemIds := make([]string, 0)
	for _, consume := range selfConsumes {
		itemIds = append(itemIds, consume.ItemId)
	}

	itemMap, err := getItemMapByIds(ctx, itemIds)
	if err != nil {
		logutil.Errorf("query item map failed, err:%v", err)
		return nil, err
	}

	results := make([]*vo.SelfConsumeVO, 0)
	for _, consume := range selfConsumes {
		consumeVO := &vo.SelfConsumeVO{
			Id:         consume.Id.Hex(),
			ItemId:     consume.ItemId,
			Count:      consume.Count,
			UpdateTime: util.FormatTime(consume.UpdateTime),
		}
		item := itemMap[consume.ItemId]
		if item != nil {
			consumeVO.ItemName = item.Name
		}
		results = append(results, consumeVO)

	}

	return results, nil
}
