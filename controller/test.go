package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestMethod(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
	return
}