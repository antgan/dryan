package main

import (
	"dryan/common"
	"dryan/controller"
	_ "dryan/db"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	testGroup := g.Group("/test")
	{
		testGroup.POST("/post1", controller.TestMethod)
	}
	g.Run(":" + common.Config.PORT)
}
