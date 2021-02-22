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

func AddItem(ctx context.Context, req *vo.AddItemReq) error {
	item := &do.Item{
		Id:             bson.NewObjectId(),
		Name:           req.Name,
		SalePrice:      req.SalePrice,
		DiscountPrice1: req.DiscountPrice1,
		DiscountPrice2: req.DiscountPrice2,
		OfficialPrice:  req.OfficialPrice,
		DirectorPrice:  req.DirectorPrice,
		CreateTime:     time.Now(),
	}
	return dao.ItemOp.Insert(ctx, item)
}

func QueryItemById(ctx context.Context, req *vo.QueryByIdReq) (*vo.ItemVO, error) {
	objId, _ := util.StringToObjectId(req.Id)
	var item *do.Item
	err := dao.ItemOp.FindById(ctx, objId, &item)
	if err != nil {
		logutil.Errorf("query item by id failed, err:%v", err)
		return nil, err
	}

	itemVO := populateItemVO(item)
	return itemVO, nil
}

func populateItemVO(item *do.Item) *vo.ItemVO {
	itemVO := &vo.ItemVO{
		Id:             item.Id.Hex(),
		Name:           item.Name,
		SalePrice:      item.SalePrice,
		DiscountPrice1: item.DiscountPrice1,
		DiscountPrice2: item.DiscountPrice2,
		OfficialPrice:  item.OfficialPrice,
		DirectorPrice:  item.DirectorPrice,
		CreateTime:     item.CreateTime,
	}
	return itemVO
}

func QueryItemByIds(ctx context.Context, req *vo.QueryByIdsReq) ([]*vo.ItemVO, error) {
	results := make([]*vo.ItemVO, 0)
	items, err := queryItemByIds(ctx, req.Ids)
	if err != nil {
		return results, err
	}
	results = populateItemVOs(items, results)
	return results, nil
}

func queryItemByIds(ctx context.Context, ids []string) ([]*do.Item, error) {
	objIds := util.StringsToObjectIds(ids)
	q := bson.M{"_id": bson.M{"$in": objIds}}
	items := make([]*do.Item, 0)
	err := dao.ItemOp.Find(ctx, &items, q, nil, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query item by ids failed, err:%v", err)
		return nil, err
	}
	return items, nil
}

func QueryAllItem(ctx context.Context) ([]*vo.ItemVO, error) {
	results := make([]*vo.ItemVO, 0)

	items := make([]*do.Item, 0)
	err := dao.ItemOp.Find(ctx, &items, nil, nil, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query item by ids failed, err:%v", err)
		return nil, err
	}

	results = populateItemVOs(items, results)
	return results, nil
}

func populateItemVOs(items []*do.Item, results []*vo.ItemVO) []*vo.ItemVO {
	for _, item := range items {
		results = append(results, populateItemVO(item))
	}
	return results
}
