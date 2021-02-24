package main

import (
	"dryan/common"
	"dryan/controller"
	_ "dryan/db"
	"dryan/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.Use(middleware.Cors())
	dryanGroup := g.Group("/dryan")
	{
		dryanGroup.POST("/user/add", controller.AddUser)
		dryanGroup.POST("/user/login", controller.QueryUserByName)

		//库存查询
		dryanGroup.POST("/stock/get-all", controller.QueryStockByUserId)
		dryanGroup.POST("/self-consume/add", controller.AddSelfConsume)
		dryanGroup.POST("/self-consume/get-all", controller.QueryAllSelfConsume)

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
