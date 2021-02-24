package controller

import (
	"dryan/model"
	"dryan/model/vo"
	"dryan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddPrePurchase(ctx *gin.Context) {
	TAG := "AddPrePurchase"
	var req vo.PrePurchase
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.AddPrePurchase(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}

func QueryAllPrePurchase(ctx *gin.Context) {
	TAG := "QueryAllPrePurchase"
	var req vo.QueryByUserIdReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	allPrePurchases, err := service.QueryAllPrePurchase(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewDataResponse(allPrePurchases, TAG))
}
