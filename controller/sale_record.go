package controller

import (
	"dryan/model"
	"dryan/model/vo"
	"dryan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddSaleRecord(ctx *gin.Context) {
	TAG := "AddSaleRecord"
	var req vo.SaleRecordVO
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.AddSaleRecord(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}

func QueryAllSaleRecordUserId(ctx *gin.Context) {
	TAG := "QueryAllSaleRecordUserId"
	var req vo.QueryByUserIdReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	allSaleRecords, err := service.QueryAllSaleRecordByUserId(ctx, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewDataResponse(allSaleRecords, TAG))
}

func UpdateCustomerInfo(ctx *gin.Context) {
	TAG := "UpdateCustomerInfo"
	var req vo.UpdateCustomerInfoReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.UpdateCustomerInfo(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}
