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

func AddPrePurchase(ctx context.Context, req *vo.PrePurchase) error {
	prePurchase, err := queryPrePurchaseByName(ctx, req.Name)
	if err != nil {
		return err
	}
	if prePurchase != nil {
		return errors.New("pre purchase already exists")
	}

	for _, item := range req.Items {
		prePurchaseItem := &do.PrePurchase{
			Id:         bson.NewObjectId(),
			UserId:     req.UserId,
			Name:       req.Name,
			ItemId:     item.ItemId,
			Count:      item.Count,
			CreateTime: time.Now(),
		}
		err := dao.PrePurchaseOp.Insert(ctx, prePurchaseItem)
		if err != nil {
			logutil.Errorf("insert pre purchase failed, err:%v", err)
			return err
		}
	}
	return nil
}

func QueryAllPrePurchase(ctx context.Context, req *vo.QueryByUserIdReq) ([]*vo.PrePurchase, error) {
	q := bson.M{"user_id": req.UserId}
	prePurchases := make([]*do.PrePurchase, 0)
	err := dao.PrePurchaseOp.Find(ctx, &prePurchases, q, []string{"-create_time"}, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query all pre purchase failed, err:%v", err)
		return nil, err
	}

	itemIds := make([]string, 0)
	for _, purchase := range prePurchases {
		if !util.Contains(itemIds, purchase.ItemId) {
			itemIds = append(itemIds, purchase.ItemId)
		}
	}

	itemMap, err := getItemMapByIds(ctx, itemIds)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string][]*vo.PrePurchaseItem)
	for _, purchase := range prePurchases {
		purchaseItem := &vo.PrePurchaseItem{
			ItemId: purchase.ItemId,
			Count:  purchase.Count,
		}
		if item, ok := itemMap[purchase.ItemId]; ok {
			purchaseItem.ItemName = item.Name
		}
		resultMap[purchase.Name] = append(resultMap[purchase.Name], purchaseItem)
	}

	results := make([]*vo.PrePurchase, 0)
	for name, items := range resultMap {
		results = append(results, &vo.PrePurchase{
			Name:  name,
			Items: items,
		})
	}

	//排序
	sortName := make([]string, 0)
	for _, prePurchase := range prePurchases {
		if !util.Contains(sortName, prePurchase.Name) {
			sortName = append(sortName, prePurchase.Name)
		}
	}
	sortResults := make([]*vo.PrePurchase, 0)
	for _, name := range sortName {
		for _, result := range results {
			if name == result.Name {
				sortResults = append(sortResults, result)
			}
		}
	}
	return sortResults, nil
}

func queryPrePurchaseByName(ctx context.Context, name string) (*do.PrePurchase, error) {
	q := bson.M{"name": name}

	prePurchases := make([]*do.PrePurchase, 0)
	err := dao.PrePurchaseOp.Find(ctx, &prePurchases, q, nil, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query pre purchase by name failed, err:%v", err)
		return nil, err
	}

	if len(prePurchases) > 0 {
		return prePurchases[0], nil
	}
	return nil, nil
}
