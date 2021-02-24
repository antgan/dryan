package controller

import (
	"dryan/model"
	"dryan/model/vo"
	"dryan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddSelfConsume(ctx *gin.Context) {
	TAG := "AddSelfConsume"
	var req vo.AddSelfConsumeReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	err := service.AddSelfConsume(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewSuccessResponse(TAG))
}

func QueryAllSelfConsume(ctx *gin.Context) {
	TAG := "QueryAllSelfConsume"
	var req vo.QueryByUserIdReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, model.NewBindFailedResponse(TAG))
		return
	}

	consumes, err := service.QueryAllSelfConsume(ctx, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, model.NewErrorResponse(err, TAG))
		return
	}

	ctx.JSON(http.StatusOK, model.NewDataResponse(consumes, TAG))
}
