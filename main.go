package main

import (
	"dryan/common"
	"dryan/controller"
	_ "dryan/db"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	dryanGroup := g.Group("/dryan")
	{
		dryanGroup.POST("/user/login", controller.QueryUserByName)

		//商品管理
		dryanGroup.POST("/item/add", controller.AddItem)
		dryanGroup.POST("/item/get-all", controller.QueryAllItem)

		//采购预设
		dryanGroup.POST("/pre-purchase/add", controller.AddPrePurchase)
		dryanGroup.POST("/pre-purchase/get-all", controller.QueryAllPrePurchase)

		//入货管理
		dryanGroup.POST("/purchase/add", controller.AddPurchaseRecord)
		dryanGroup.POST("/purchase/get-all", controller.QueryAllPurchaseByUserId)

		//销售管理
		dryanGroup.POST("/sale/add", controller.AddSaleRecord)
		dryanGroup.POST("/sale/get-all", controller.QueryAllSaleRecordUserId)
		dryanGroup.POST("/sale/update-customer-info", controller.UpdateCustomerInfo)
	}

	g.Run(":" + common.Config.PORT)
}
