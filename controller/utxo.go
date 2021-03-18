package controller

import (
	"encoding/hex"
	"log"
	"net/http"
	"satoblock/lib/utils"
	"satoblock/model"
	"satoblock/service"

	"github.com/gin-gonic/gin"
)

// GetUtxoByAddress
// @Summary 通过地址address获取相关utxo列表
// @Tags UTXO
// @Produce  json
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /address/{address}/utxo [get]
func GetUtxoByAddress(ctx *gin.Context) {
	log.Printf("GetUtxoByAddress enter")

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		log.Printf("address invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, err := service.GetUtxoByAddress(addressPkh)
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetUtxoByGenesis
// @Summary 通过溯源genesis获取相关utxo列表
// @Tags UTXO
// @Produce  json
// @Param genesis path string true "Genesis" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.TxOutResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /genesis/{genesis}/utxo [get]
func GetUtxoByGenesis(ctx *gin.Context) {
	log.Printf("GetUtxoByGenesis enter")

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		log.Printf("genesisId invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetUtxoByGenesis(hex.EncodeToString(genesisId))
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
