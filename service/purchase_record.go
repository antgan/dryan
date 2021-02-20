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
	"gopkg.in/mgo.v2/bson"
	"time"
)

func AddPurchaseRecord(ctx context.Context, purchase *vo.Purchase) error {
	if len(purchase.Items) <= 0 {
		return errors.New("empty items")
	}
	user, err := QueryUserById(ctx, purchase.UserId)
	if err != nil || user == nil {
		logutil.Errorf("user not exists, err:%v", err)
		return err
	}

	//获取进货单价
	itemIds := make([]string, 0)
	for _, item := range purchase.Items {
		itemIds = append(itemIds, item.ItemId)
	}
	itemMap, err := getItemMapByIds(ctx, itemIds)
	if err != nil {
		return err
	}
	purchasePriceMapping, err := getPurchasePriceMapping(ctx, itemMap, user.Type)
	if err != nil {
		return err
	}

	purchaseRecords := make([]*do.PurchaseRecord, 0)
	serialId := util.NewUUID()
	now := time.Now()
	for _, item := range purchase.Items {
		purchaseRecord := &do.PurchaseRecord{
			Id:         bson.NewObjectId(),
			UserId:     purchase.UserId,
			UserType:   user.Type,
			SerialId:   serialId,
			ItemId:     item.ItemId,
			Count:      item.Count,
			Price:      purchasePriceMapping[item.ItemId],
			CreateTime: now,
		}

		purchaseRecords = append(purchaseRecords, purchaseRecord)
	}

	//插入入库记录和库存
	for _, record := range purchaseRecords {
		err := dao.PurchaseRecordOp.Insert(ctx, record)
		if err != nil {
			logutil.Errorf("insert purchase failed, err:%v", err)
			return err
		}
		err = UpdateStockCount(ctx, record.UserId, record.ItemId, record.Count, true)
		if err != nil {
			logutil.Errorf("update stock count failed, err:%v", err)
			return err
		}
	}

	return nil
}

func getPurchasePriceMapping(ctx context.Context, itemMap map[string]*do.Item, userType string) (map[string]int, error) {
	purchaseMap := make(map[string]int)

	//计算入货单价价格
	for _, itemDO := range itemMap {
		if userType == constant.DRYAN_USER_TYPE_OFFICIAL {
			purchaseMap[itemDO.Id.Hex()] = itemDO.OfficialPrice
		}
		if userType == constant.DRYAN_USER_TYPE_DIRECTOR {
			purchaseMap[itemDO.Id.Hex()] = itemDO.DirectorPrice
		}
	}

	return purchaseMap, nil
}

func getSalePriceMapping(ctx context.Context, itemMap map[string]*do.Item, discountLevel int) (map[string]int, error) {
	saleMap := make(map[string]int)

	//按照折扣算出销售价格
	for _, itemDO := range itemMap {
		if discountLevel == 1 {
			saleMap[itemDO.Id.Hex()] = itemDO.DiscountPrice1
		} else if discountLevel == 2 {
			saleMap[itemDO.Id.Hex()] = itemDO.DiscountPrice2
		} else {
			saleMap[itemDO.Id.Hex()] = itemDO.SalePrice
		}
	}

	return saleMap, nil
}

func QueryPurchaseByUserId(ctx context.Context, userId string) ([]*vo.Purchase, error) {
	purchaseRecords := make([]*do.PurchaseRecord, 0)
	q := bson.M{"user_id": userId}
	sort := []string{"-create_time"}
	err := dao.PurchaseRecordOp.Find(ctx, &purchaseRecords, q, sort, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query purchase record failed, err:%v", err)
		return nil, err
	}

	//收集所有商品id，为了补充name
	allItemIds := make([]string, 0)
	for _, record := range purchaseRecords {
		if !util.Contains(allItemIds, record.ItemId) {
			allItemIds = append(allItemIds, record.ItemId)
		}
	}
	itemNameMapping, err := getItemNameMapping(ctx, allItemIds)
	if err != nil {
		return nil, err
	}

	//按照流水号分组排序
	itemsGroupBySerial := make(map[string][]*vo.PurchaseItem)
	timeGroupBySerial := make(map[string]time.Time)
	for _, record := range purchaseRecords {
		itemsGroupBySerial[record.SerialId] = append(itemsGroupBySerial[record.SerialId], &vo.PurchaseItem{
			ItemId:        record.ItemId,
			ItemName:      itemNameMapping[record.ItemId],
			Count:         record.Count,
			PurchasePrice: record.Price,
		})
		timeGroupBySerial[record.SerialId] = record.CreateTime
	}

	//聚合最终展示的vo
	results := make([]*vo.Purchase, 0)
	for serialId, items := range itemsGroupBySerial {
		results = append(results, &vo.Purchase{
			UserId:             userId,
			SerialId:           serialId,
			Items:              items,
			TotalPurchasePrice: calcTotalPriceByItems(items),
			CreateTime:         timeGroupBySerial[serialId],
		})
	}

	return results, nil
}

func getItemNameMapping(ctx context.Context, allItemIds []string) (map[string]string, error) {
	resultMap := make(map[string]string)
	itemDOs, err := queryItemByIds(ctx, allItemIds)
	if err != nil {
		return resultMap, err
	}
	for _, itemDO := range itemDOs {
		resultMap[itemDO.Id.Hex()] = itemDO.Name
	}
	return resultMap, nil
}

func calcTotalPriceByItems(items []*vo.PurchaseItem) int {
	totalPrice := 0
	for _, item := range items {
		totalPrice += item.PurchasePrice * item.Count
	}
	return totalPrice
}
