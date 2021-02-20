package vo

import "time"

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
