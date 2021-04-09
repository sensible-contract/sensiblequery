package controller

import (
	"log"
	"net/http"
	"satoblock/model"
	"satoblock/service"

	"github.com/gin-gonic/gin"
)

// Satotx
// @Summary Welcome message
// @Produce  json
// @Success 200 {object} model.Response{data=model.Welcome} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router / [get]
func Satotx(ctx *gin.Context) {
	log.Printf("Satotx enter")

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "Welcome to use sensible contract on Bitcoin SV!",
		Data: &model.Welcome{
			Contact: "",
			Job:     "",
			Github:  "https://github.com/sensible-contract",
		},
	})
}

// GetBlockchainInfo 获取最新区块位置、同步状态等信息
// @Summary 获取最新区块位置、同步状态等信息
// @Produce  json
// @Success 200 {object} model.Response{data=model.BlockchainInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /blockchain/info [get]
func GetBlockchainInfo(ctx *gin.Context) {
	log.Printf("GetBlockchainInfo enter")

	blk, err := service.GetBestBlock()
	if err != nil {
		log.Printf("best block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get best block failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: &model.BlockchainInfoResp{
			Chain:         "main",
			Blocks:        blk.Height + 1,
			Headers:       blk.Height + 1,
			BestBlockHash: blk.BlockIdHex,
			Difficulty:    "",
			MedianTime:    blk.BlockTime,
			Chainwork:     "",
		},
	})

}
