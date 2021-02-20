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
		dryanGroup.POST("/item/add", controller.AddItem)
		dryanGroup.POST("/item/get-all", controller.QueryAllItem)
	}

	g.Run(":" + common.Config.PORT)
}
