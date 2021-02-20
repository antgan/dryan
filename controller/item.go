package controller

import (
	"dryan/model"
	"dryan/model/vo"
	"dryan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddItem(ctx *gin.Context) {
	TAG := "AddItem"
	var req vo.AddItemReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.AddItem(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}

func QueryAllItem(ctx *gin.Context) {
	TAG := "QueryAllItem"

	allItems, err := service.QueryAllItem(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewDataResponse(allItems, TAG))
}
