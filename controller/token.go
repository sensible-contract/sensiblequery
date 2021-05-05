package controller

import (
	"log"
	"net/http"
	"satosensible/model"
	"satosensible/service"

	"github.com/gin-gonic/gin"
)

// ListAllTokenInfo
// @Summary 查询Token简述
// @Tags token
// @Produce  json
// @Success 200 {object} model.Response{data=model.TokenInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /token/info [get]
func ListAllTokenInfo(ctx *gin.Context) {
	log.Printf("ListAllTokenInfo enter")

	result, err := service.GetTokenInfo()
	if err != nil {
		log.Printf("get dummy failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get dummy failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
