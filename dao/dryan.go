package dao

import (
	"dryan/constant"
	"dryan/db/mongo"
)

var ItemOp *mongo.MongoOp
var StockOp *mongo.MongoOp
var PurchaseRecordOp *mongo.MongoOp
var SaleRecordOp *mongo.MongoOp
var SaleRecordSummaryOp *mongo.MongoOp
var PrePurchaseOp *mongo.MongoOp
var UserOp *mongo.MongoOp

func init() {
	ItemOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_ITEM)
	StockOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_STOCK)
	PurchaseRecordOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_PURCHASE_RECORD)
	SaleRecordOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_SALE_RECORD)
	SaleRecordSummaryOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_SALE_RECORD_SUMMARY)
	PrePurchaseOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_PRE_PURCHASE)
	UserOp = mongo.NewMongoOp(constant.MONGO_DB_DRYAN, constant.MONGO_C_DRYAN_USER)
}
