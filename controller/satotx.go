package controller

import (
	"log"
	"net/http"
	"satoblock/model"

	"github.com/gin-gonic/gin"
)

// Satotx
func Satotx(ctx *gin.Context) {
	log.Printf("Satotx enter")

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "Welcome to use sensible contract on Bitcoin SV!",
		Data: &model.Welcome{
			Contact: "",
			Job:     "",
			Github:  "https://github.com/sensible-group",
		},
	})
}

func GetBlockchainInfo(ctx *gin.Context) {
	log.Printf("GetBlockchainInfo enter")

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: &model.BlockchainInfoResp{
			Chain:         "main",
			Blocks:        0,
			Headers:       0,
			BestBlockHash: "",
			Difficulty:    "",
			MedianTime:    0,
			Chainwork:     "",
		},
	})

}
