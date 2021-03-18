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

// GetHistoryByAddress
// @Summary 通过地址address获取相关tx历史列表，返回详细输入/输出
// @Tags History
// @Produce  json
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /address/{address}/history [get]
func GetHistoryByAddress(ctx *gin.Context) {
	log.Printf("GetHistoryByAddress enter")

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		log.Printf("address invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, err := service.GetHistoryByAddress(hex.EncodeToString(addressPkh))
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

// GetHistoryByGenesis
// @Summary 通过溯源genesis获取相关tx历史列表，返回详细输入/输出
// @Tags History
// @Produce  json
// @Param genesis path string true "Genesis" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /genesis/{genesis}/history [get]
func GetHistoryByGenesis(ctx *gin.Context) {
	log.Printf("GetHistoryByGenesis enter")

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		log.Printf("genesisId invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetHistoryByGenesis(hex.EncodeToString(genesisId))
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

// GetFTHistoryByGenesis
// @Summary 通过FT合约CodeHash+溯源genesis获取地址相关tx历史列表，返回详细输入/输出
// @Tags History, token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/history/{codehash}/{genesis}/{address} [get]
func GetFTHistoryByGenesis(ctx *gin.Context) {
	log.Printf("GetFTHistoryByGenesis enter")

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		log.Printf("genesisId invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetHistoryByGenesis(hex.EncodeToString(genesisId))
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

// GetNFTHistoryByGenesis
// @Summary 通过FT合约CodeHash+溯源genesis获取地址相关tx历史列表，返回详细输入/输出
// @Tags History, token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/history/{codehash}/{genesis}/{address} [get]
func GetNFTHistoryByGenesis(ctx *gin.Context) {
	log.Printf("GetNFTHistoryByGenesis enter")

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		log.Printf("genesisId invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetHistoryByGenesis(hex.EncodeToString(genesisId))
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
