package main

import (
	"dryan/common"
	"dryan/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	testGroup := g.Group("/test")
	{
		testGroup.POST("/post1", controller.TestMethod)
	}
	g.Run(common.Config.PORT)
}

func init() {
	common.InitConf()
}
