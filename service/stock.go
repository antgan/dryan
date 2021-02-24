package service

import (
	"context"
	"dryan/dao"
	"dryan/model/do"
	"dryan/model/vo"
	"dryan/util"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func UpdateStockCount(ctx context.Context, userId, itemId string, count int, isAdd bool) error {
	incCount := 0
	if isAdd {
		incCount = count
	} else {
		incCount = count * -1
	}
	q := bson.M{"user_id": userId, "item_id": itemId}
	updateMap := bson.M{"$inc": bson.M{"remain_count": incCount}, "$set": bson.M{"update_time": time.Now()}}
	err := dao.StockOp.Upsert(ctx, q, updateMap)
	if err != nil {
		logutil.Errorf("update stock count failed, err:%v", err)
		return err
	}

	return nil
}

func QueryStockByUserId(ctx context.Context, userId string) ([]*vo.StockVO, error) {
	stocks, err := queryStockByUserId(ctx, userId)
	if err != nil {
		logutil.Errorf("query stock by user id failed, userId:%s, err:%v", userId, err)
		return nil, err
	}

	itemIds := make([]string, 0)
	for _, stock := range stocks {
		itemIds = append(itemIds, stock.ItemId)
	}

	itemMap, err := getItemMapByIds(ctx, itemIds)
	if err != nil {
		logutil.Errorf("query item map by ids failed, err:%v", err)
		return nil, err
	}

	results := make([]*vo.StockVO, 0)
	for _, stock := range stocks {
		stockVO := &vo.StockVO{
			Id:          stock.Id.Hex(),
			UserId:      stock.UserId,
			ItemId:      stock.ItemId,
			RemainCount: stock.RemainCount,
			UpdateTime:  util.FormatTime(stock.UpdateTime),
		}
		itemDO := itemMap[stock.ItemId]
		if itemDO != nil {
			stockVO.ItemName = itemDO.Name
		}
		results = append(results, stockVO)
	}
	return results, nil
}

func queryStockByUserId(ctx context.Context, userId string) ([]*do.Stock, error) {
	q := bson.M{"user_id": userId}
	stockDOs := make([]*do.Stock, 0)
	err := dao.StockOp.Find(ctx, &stockDOs, q, nil, nil, 0, 0)
	return stockDOs, err
}

func queryStockMapByUserId(ctx context.Context, userId string) (map[string]*do.Stock, error) {
	resultMap := make(map[string]*do.Stock)
	stocks, err := queryStockByUserId(ctx, userId)
	if err != nil {
		return resultMap, err
	}
	for _, stock := range stocks {
		resultMap[stock.ItemId] = stock
	}
	return resultMap, nil
}
