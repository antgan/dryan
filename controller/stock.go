package controller

import (
	"dryan/model"
	"dryan/model/vo"
	"dryan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func QueryStockByUserId(ctx *gin.Context) {
	TAG := "QueryStockByUserId"
	var req vo.QueryByUserIdReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	stocks, err := service.QueryStockByUserId(ctx, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewDataResponse(stocks, TAG))
}
