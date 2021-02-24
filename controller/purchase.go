package controller

import (
	"dryan/model"
	"dryan/model/vo"
	"dryan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddPurchaseRecord(ctx *gin.Context) {
	TAG := "AddPurchaseRecord"
	var req vo.Purchase
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.AddPurchaseRecord(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}

func QueryAllPurchaseByUserId(ctx *gin.Context) {
	TAG := "QueryAllPurchaseByUserId"
	var req vo.QueryByUserIdReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	allPurchases, err := service.QueryPurchaseByUserId(ctx, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewDataResponse(allPurchases, TAG))
}

func DeletePurchaseRecord(ctx *gin.Context) {
	TAG := "DeletePurchaseRecord"
	var req vo.DeleteRecordReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.DeletePurchaseRecord(ctx, req.UserId, req.SerialId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}
