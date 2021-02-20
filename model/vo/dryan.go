package vo

import (
	"time"
)

type AddItemReq struct {
	Name           string `json:"name"`
	SalePrice      int    `json:"salePrice"`      //零售价
	DiscountPrice1 int    `json:"discountPrice1"` //优惠1
	DiscountPrice2 int    `json:"discountPrice2"` //优惠2
	OfficialPrice  int    `json:"officialPrice"`  //官方价格
	DirectorPrice  int    `json:"directorPrice"`  // 董事价格
}

type QueryByIdReq struct {
	Id string `json:"id"`
}

type QueryByIdsReq struct {
	Ids []string `json:"ids"`
}

type ItemVO struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	SalePrice      int       `json:"salePrice"`      //零售价
	DiscountPrice1 int       `json:"discountPrice1"` //优惠1
	DiscountPrice2 int       `json:"discountPrice2"` //优惠2
	OfficialPrice  int       `json:"officialPrice"`  //官方价格
	DirectorPrice  int       `json:"directorPrice"`  // 董事价格
	CreateTime     time.Time `json:"createTime"`
}

type PrePurchase struct {
	Name  string             `json:"name" binding:"required"` //预设名称，唯一
	Items []*PrePurchaseItem `json:"items"`
}

type PrePurchaseItem struct {
	ItemId string `json:"itemId"`
	Count  int    `json:"count"`
}

type Purchase struct {
	UserId     string          `json:"userId"`
	SerialId   string          `json:"serialId"` //流水号
	Items      []*PurchaseItem `json:"items"`
	TotalPrice int             `json:"totalPrice"`
	CreateTime time.Time       `json:"createTime"`
}

type PurchaseItem struct {
	ItemId   string `json:"itemId"`
	ItemName string `json:"itemName"`
	Count    int    `json:"count"`
	Price    int    `json:"price"` //进货单价
}

type QueryByUserIdReq struct {
	UserId string `json:"userId" binding:"required"`
}
