package do

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

//商品
type Item struct {
	Id             bson.ObjectId `bson:"_id"`
	Name           string        `bson:"name"`
	SalePrice      int           `bson:"sale_price"`       //零售价
	DiscountPrice1 int           `bson:"discount_price_1"` //优惠1
	DiscountPrice2 int           `bson:"discount_price_2"` //优惠2
	OfficialPrice  int           `bson:"official_price"`   //官方价格
	DirectorPrice  int           `bson:"director_price"`   // 董事价格
	CreateTime     time.Time     `bson:"create_time"`
}

//库存
type Stock struct {
	Id          bson.ObjectId `bson:"_id"`
	UserId      string        `bson:"user_id"`
	ItemId      string        `bson:"item_id"`
	RemainCount int           `bson:"remain_count"`
	UpdateTime  time.Time     `bson:"update_time"`
}

//入货记录
type PurchaseRecord struct {
	Id         bson.ObjectId `bson:"_id"`
	UserId     string        `bson:"user_id"`
	UserType   string        `bson:"user_type"`
	SerialId   string        `bson:"serial_id"` //流水号
	ItemId     string        `bson:"item_id"`
	Count      int           `bson:"count"`
	Price      int           `bson:"price"` //单价
	CreateTime time.Time     `bson:"create_time"`
}

//销售记录
type SaleRecord struct {
	Id     bson.ObjectId `bson:"_id"`
	UserId string        `bson:"user_id"`

	SerialId      string    `bson:"serial_id"` //流水号
	ItemId        string    `bson:"item_id"`
	Count         int       `bson:"count"`
	PurchasePrice int       `bson:"purchase_price"`
	SalePrice     int       `bson:"sale_price"`
	Profit        int       `bson:"profit"` //利润
	CreateTime    time.Time `bson:"create_time"`
	UpdateTime    time.Time `bson:"update_time"`

	//顾客信息
	CustomerName  string    `bson:"customer_name"`
	Address       string    `bson:"address"`
	Logistics     string    `bson:"logistics"`
	ExpressNumber string    `bson:"express_number"`
	ExpressTime   time.Time `bson:"express_time"`
}

//预设进货套餐
type PrePurchase struct {
	Id         bson.ObjectId `bson:"_id"`
	Name       string        `bson:"name"` //预设名称，唯一
	ItemId     string        `bson:"item_id"`
	Count      int           `bson:"count"`
	CreateTime time.Time     `bson:"create_time"`
}
