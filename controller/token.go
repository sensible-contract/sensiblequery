package controller

import (
	"net/http"
	"sensiblequery/logger"
	"sensiblequery/model"
	"sensiblequery/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListAllTokenInfo
// @Summary 查询Token简述
// @Tags token
// @Produce  json
// @Success 200 {object} model.Response{data=model.TokenInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /token/info [get]
func ListAllTokenInfo(ctx *gin.Context) {
	logger.Log.Info("ListAllTokenInfo enter")

	result, err := service.GetTokenInfo()
	if err != nil {
		logger.Log.Info("get dummy failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get dummy failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
