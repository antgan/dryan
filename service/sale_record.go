package service

import (
	"context"
	"dryan/dao"
	"dryan/model/do"
	"dryan/model/vo"
	"dryan/util"
	"errors"
	"fmt"
	logutil "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func AddSaleRecord(ctx context.Context, saleRecord *vo.SaleRecordVO) error {
	if len(saleRecord.Items) <= 0 {
		return errors.New("empty items")
	}
	user, err := QueryUserById(ctx, saleRecord.UserId)
	if err != nil || user == nil {
		logutil.Errorf("user not exists, err:%v", err)
		return err
	}

	//获取进货单价
	itemIds := make([]string, 0)
	for _, item := range saleRecord.Items {
		itemIds = append(itemIds, item.ItemId)
	}
	itemMap, err := getItemMapByIds(ctx, itemIds)
	if err != nil {
		return err
	}
	//校验库存还够不够
	err = checkStockByItem(ctx, saleRecord, itemMap)
	if err != nil {
		return err
	}

	purchasePriceMapping, err := getPurchasePriceMapping(ctx, itemMap, user.Type)
	if err != nil {
		return err
	}

	//获取销售单价
	discountLevel := getDiscountLevel(saleRecord, itemMap)
	salePriceMapping, err := getSalePriceMapping(ctx, itemMap, discountLevel)
	if err != nil {
		return err
	}

	saleRecords := make([]*do.SaleRecord, 0)
	serialId := util.NewUUID()
	now := time.Now()
	totalProfit := 0
	totalSalePrice := 0
	totalPurchasePrice := 0
	for _, item := range saleRecord.Items {
		profit := (salePriceMapping[item.ItemId] - purchasePriceMapping[item.ItemId]) * item.Count
		saleRecordDO := &do.SaleRecord{
			Id:            bson.NewObjectId(),
			UserId:        saleRecord.UserId,
			SerialId:      serialId,
			ItemId:        item.ItemId,
			Count:         item.Count,
			PurchasePrice: purchasePriceMapping[item.ItemId],
			SalePrice:     salePriceMapping[item.ItemId],
			CreateTime:    now,
			UpdateTime:    now,
			Profit:        profit,
		}

		saleRecords = append(saleRecords, saleRecordDO)
		totalProfit += profit
		totalSalePrice += saleRecordDO.SalePrice * saleRecordDO.Count
		totalPurchasePrice += saleRecordDO.PurchasePrice * saleRecordDO.Count
	}

	//插入销售记录和库存
	for _, record := range saleRecords {
		err = dao.SaleRecordOp.Insert(ctx, record)
		if err != nil {
			logutil.Errorf("insert sale record failed, err:%v", err)
			return err
		}
		err = UpdateStockCount(ctx, record.UserId, record.ItemId, record.Count, false)
		if err != nil {
			logutil.Errorf("update stock count failed, err:%v", err)
			return err
		}
	}
	saleSummary := &do.SaleRecordSummary{
		Id:                 bson.NewObjectId(),
		UserId:             saleRecord.UserId,
		SerialId:           serialId,
		Profit:             totalProfit,
		TotalPurchasePrice: totalPurchasePrice,
		TotalSalePrice:     totalSalePrice,
		CustomerName:       saleRecord.CustomerName,
		Address:            saleRecord.Address,
		CreateTime:         now,
		UpdateTime:         now,
	}

	err = dao.SaleRecordSummaryOp.Insert(ctx, saleSummary)
	if err != nil {
		logutil.Errorf("insert sale summary failed, err:%v", err)
		return err
	}

	return nil
}

func checkStockByItem(ctx context.Context, saleRecord *vo.SaleRecordVO, itemMap map[string]*do.Item) error {
	stockMap, err := queryStockMapByUserId(ctx, saleRecord.UserId)
	if err != nil {
		logutil.Errorf("query stock failed, err:%v", err)
		return err
	}
	for _, item := range saleRecord.Items {
		stock := stockMap[item.ItemId]
		if item.Count > stock.RemainCount {
			return errors.New(fmt.Sprintf("%s库存不足", itemMap[item.ItemId].Name))
		}
	}
	return nil
}

func getDiscountLevel(saleRecord *vo.SaleRecordVO, itemMap map[string]*do.Item) int {
	fullPrice := 0 //按照零售价总和，看是否达到优惠门槛
	for _, saleItem := range saleRecord.Items {
		fullPrice += itemMap[saleItem.ItemId].SalePrice * saleItem.Count
	}
	isMixed := len(saleRecord.Items) > 1 //混装
	discountLevel := 0
	if isMixed {
		if fullPrice < 1850 {
			discountLevel = 0
		}
		if fullPrice >= 1850 && fullPrice < 2600 {
			discountLevel = 1
		}
		if fullPrice >= 2600 {
			discountLevel = 2
		}
	} else {
		if fullPrice >= 650 && fullPrice < 1850 {
			discountLevel = 1
		}
		if fullPrice >= 2600 {
			discountLevel = 2
		}
	}

	return discountLevel
}

func getItemMapByIds(ctx context.Context, itemIds []string) (map[string]*do.Item, error) {
	itemMap := make(map[string]*do.Item)
	itemDOs, err := queryItemByIds(ctx, itemIds)
	if err != nil {
		return nil, err
	}
	for _, itemDO := range itemDOs {
		itemMap[itemDO.Id.Hex()] = itemDO
	}
	return itemMap, nil
}

func QueryAllSaleRecordByUserId(ctx context.Context, userId string) ([]*vo.SaleRecordVO, error) {
	saleRecords := make([]*do.SaleRecord, 0)
	q := bson.M{"user_id": userId}
	sort := []string{"-create_time"}
	err := dao.SaleRecordOp.Find(ctx, &saleRecords, q, sort, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query sale record failed, err:%v", err)
		return nil, err
	}

	//排序用
	sortSerialIds := make([]string, 0)
	for _, record := range saleRecords {
		if !util.Contains(sortSerialIds, record.SerialId) {
			sortSerialIds = append(sortSerialIds, record.SerialId)
		}
	}

	//收集所有商品id，为了补充name
	allItemIds := make([]string, 0)
	for _, record := range saleRecords {
		if !util.Contains(allItemIds, record.ItemId) {
			allItemIds = append(allItemIds, record.ItemId)
		}
	}
	itemNameMapping, err := getItemNameMapping(ctx, allItemIds)
	if err != nil {
		return nil, err
	}

	//按照流水号分组排序
	itemsGroupBySerial := make(map[string][]*vo.SaleItem)
	for _, record := range saleRecords {
		itemsGroupBySerial[record.SerialId] = append(itemsGroupBySerial[record.SerialId], &vo.SaleItem{
			ItemId:        record.ItemId,
			ItemName:      itemNameMapping[record.ItemId],
			Count:         record.Count,
			PurchasePrice: record.PurchasePrice,
			SalePrice:     record.SalePrice,
		})
	}

	//聚合最终展示的vo
	summaryMap, err := getSaleSummaryMapBySerialIds(ctx, userId, sortSerialIds)
	if err != nil {
		logutil.Errorf("get sale summary by serial ids failed, err:%v", err)
		return nil, err
	}

	results := make([]*vo.SaleRecordVO, 0)
	for serialId, items := range itemsGroupBySerial {
		summary := summaryMap[serialId]
		if summary == nil {
			logutil.Errorf("summary not found, serialId:%s", serialId)
			return nil, errors.New("summary not found")
		}
		results = append(results, &vo.SaleRecordVO{
			UserId:             userId,
			SerialId:           serialId,
			Items:              items,
			Profit:             summary.Profit,
			TotalSalePrice:     summary.TotalSalePrice,
			TotalPurchasePrice: summary.TotalPurchasePrice,
			CustomerName:       summary.CustomerName,
			Address:            summary.Address,
			Logistics:          summary.Logistics,
			ExpressNumber:      summary.ExpressNumber,
			ExpressTime:        util.FormatTime(summary.ExpressTime),
			CreateTime:         util.FormatTime(summary.CreateTime),
			UpdateTime:         util.FormatTime(summary.UpdateTime),
		})
	}

	//Map无序，按照DB顺序排序
	sortResults := make([]*vo.SaleRecordVO, 0)
	for _, sortSerialId := range sortSerialIds {
		for _, result := range results {
			if sortSerialId == result.SerialId {
				sortResults = append(sortResults, result)
			}
		}
	}

	return sortResults, nil
}

func getSaleSummaryBySerialId(ctx context.Context, userId, serialId string) (*do.SaleRecordSummary, error) {
	summaries := make([]*do.SaleRecordSummary, 0)
	q := bson.M{"user_id": userId, "serial_id": serialId}
	err := dao.SaleRecordSummaryOp.Find(ctx, &summaries, q, nil, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query sale sumary failed, err:%v", err)
		return nil, err
	}

	if len(summaries) > 0 {
		return summaries[0], nil
	}
	return nil, errors.New("not single summary")
}

func getSaleSummaryMapBySerialIds(ctx context.Context, userId string, serialIds []string) (map[string]*do.SaleRecordSummary, error) {
	resultMap := make(map[string]*do.SaleRecordSummary)
	summaries := make([]*do.SaleRecordSummary, 0)
	q := bson.M{"user_id": userId, "serial_id": bson.M{"$in": serialIds}}
	err := dao.SaleRecordSummaryOp.Find(ctx, &summaries, q, nil, nil, 0, 0)
	if err != nil {
		logutil.Errorf("query sale sumary failed, err:%v", err)
		return resultMap, err
	}

	for _, summary := range summaries {
		resultMap[summary.SerialId] = summary
	}
	return resultMap, nil
}

func UpdateCustomerInfo(ctx context.Context, req *vo.UpdateCustomerInfoReq) error {
	summary, err := getSaleSummaryBySerialId(ctx, req.UserId, req.SerialId)
	if err != nil {
		logutil.Errorf("query sale summary failed, err:%v", err)
		return err
	}
	if summary == nil {
		return errors.New("empty summary")
	}
	updateMap := bson.M{"$set": bson.M{
		"customer_name":  req.CustomerName,
		"address":        req.Address,
		"logistics":      req.Logistics,
		"express_number": req.ExpressNumber,
		"express_time":   time.Now(),
	}}

	err = dao.SaleRecordSummaryOp.UpdateById(ctx, summary.Id, updateMap)
	if err != nil {
		logutil.Errorf("update summary failed, err:%v", err)
		return err
	}

	return nil
}

func DeleteSaleRecord(ctx context.Context, userId string, serialId string) error {
	q := bson.M{"user_id": userId, "serial_id": serialId}
	saleRecords := make([]*do.SaleRecord, 0)
	err := dao.SaleRecordOp.Find(ctx, &saleRecords, q, nil, nil, 0, 0)
	if err != nil {
		logutil.Errorf("find sale record by serial failed, err:%v", err)
		return err
	}

	for _, record := range saleRecords {
		//恢复库存
		err = UpdateStockCount(ctx, userId, record.ItemId, record.Count, true)
		if err != nil {
			logutil.Errorf("update stock for delete sale record failed, err:%v", err)
			return err
		}
		err = dao.SaleRecordOp.DeleteById(ctx, record.Id)
		if err != nil {
			logutil.Errorf("delete sale record failed, err:%v", err)
			return err
		}
	}
	//删除summary
	err = dao.SaleRecordSummaryOp.Delete(ctx, q)
	if err != nil {
		logutil.Errorf("delete summary failed, serialId:%s, err:%v", serialId, err)
		return err
	}

	return nil
}
