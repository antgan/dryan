package vo

type AddUserReq struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type"` //official官方；director董事
}

type UserVO struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"` //official官方；director董事
	CreateTime string `json:"createTime"`
}

type AddItemReq struct {
	Name           string `json:"name" binding:"required"`
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
	Id             string `json:"id"`
	Name           string `json:"name"`
	SalePrice      int    `json:"salePrice"`      //零售价
	DiscountPrice1 int    `json:"discountPrice1"` //优惠1
	DiscountPrice2 int    `json:"discountPrice2"` //优惠2
	OfficialPrice  int    `json:"officialPrice"`  //官方价格
	DirectorPrice  int    `json:"directorPrice"`  // 董事价格
	CreateTime     string `json:"createTime"`
}

type PrePurchase struct {
	UserId string             `json:"userId" binding:"required"`
	Name   string             `json:"name" binding:"required"` //预设名称，唯一
	Items  []*PrePurchaseItem `json:"items"`
}

type PrePurchaseItem struct {
	ItemId   string `json:"itemId"`
	ItemName string `json:"itemName"`
	Count    int    `json:"count"`
}

type Purchase struct {
	UserId             string          `json:"userId"`
	SerialId           string          `json:"serialId"` //流水号
	Items              []*PurchaseItem `json:"items"`
	TotalPurchasePrice int             `json:"totalPurchasePrice"`
	CreateTime         string          `json:"createTime"`
}

type PurchaseItem struct {
	ItemId        string `json:"itemId"`
	ItemName      string `json:"itemName"`
	Count         int    `json:"count"`
	PurchasePrice int    `json:"purchasePrice"` //进货单价
}

type QueryByUserIdReq struct {
	UserId string `json:"userId" binding:"required"`
}

type QueryUserByNameReq struct {
	Name string `json:"name" binding:"required"`
}

//销售记录
type SaleRecordVO struct {
	UserId   string `json:"userId"`
	SerialId string `json:"serialId"` //流水号

	Items []*SaleItem `json:"items"`

	Profit             int    `json:"profit"`             //利润
	TotalSalePrice     int    `json:"totalSalePrice"`     //销售价格
	TotalPurchasePrice int    `json:"totalPurchasePrice"` //入货价格
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`

	//顾客信息
	CustomerName  string `json:"customerName"`
	Address       string `json:"address"`
	Logistics     string `json:"logistics"`
	ExpressNumber string `json:"expressNumber"`
	ExpressTime   string `json:"expressTime"`
}

type SaleItem struct {
	ItemId        string `json:"itemId"`
	ItemName      string `json:"itemName"`
	Count         int    `json:"count"`
	PurchasePrice int    `json:"purchasePrice"`
	SalePrice     int    `json:"salePrice"`
}

type UpdateCustomerInfoReq struct {
	UserId   string `json:"userId"`
	SerialId string `json:"serialId"` //流水号
	//顾客信息
	CustomerName  string `json:"customerName"`
	Address       string `json:"address"`
	Logistics     string `json:"logistics"`
	ExpressNumber string `json:"expressNumber"`
}

//库存
type StockVO struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	ItemId      string `json:"itemId"`
	ItemName    string `json:"itemName"`
	RemainCount int    `json:"remainCount"`
	UpdateTime  string `json:"updateTime"`
}

type AddSelfConsumeReq struct {
	UserId string      `json:"userId"`
	Items  []*SaleItem `json:"items"`
}

//偷吃记录
type SelfConsumeVO struct {
	Id         string `json:"id"`
	ItemId     string `json:"itemId"`
	ItemName   string `json:"itemName"`
	Count      int    `json:"count"`
	UpdateTime string `json:"updateTime"`
}

type DeleteRecordReq struct {
	UserId   string `json:"userId" binding:"required"`
	SerialId string `json:"serialId" binding:"required` //流水号
}
