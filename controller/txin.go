package controller

import (
	"encoding/hex"
	"log"
	"net/http"
	"satoblock/lib/utils"
	"satoblock/model"
	"satoblock/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTxInputsByTxId
// @Summary 通过交易txid获取交易所有输入信息列表
// @Tags Txin
// @Produce  json
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Success 200 {object} model.Response{data=[]model.TxInResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /tx/{txid}/ins [get]
func GetTxInputsByTxId(ctx *gin.Context) {
	log.Printf("GetTxInputsByTxId enter")

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	result, err := service.GetTxInputsByTxId(hex.EncodeToString(txId))
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txin failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxInputsByTxIdInsideHeight
// @Summary 通过交易txid和交易被打包的区块高度height获取交易所有输入信息列表
// @Tags Txin
// @Produce  json
// @Param height path int true "Block Height" default(170)
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Success 200 {object} model.Response{data=[]model.TxInResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /height/{height}/tx/{txid}/ins [get]
func GetTxInputsByTxIdInsideHeight(ctx *gin.Context) {
	log.Printf("GetTxInputsByTxIdInsideHeight enter")

	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil {
		log.Printf("height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "height invalid"})
		return
	}

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	result, err := service.GetTxInputsByTxIdInsideHeight(blkHeight, hex.EncodeToString(txId))
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txin failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxInputByTxIdAndIdx
// @Summary 通过交易txid和index获取指定交易输入信息
// @Tags Txin
// @Produce  json
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Param index path int true "input index" default(0)
// @Success 200 {object} model.Response{data=model.TxInResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /tx/{txid}/in/{index} [get]
func GetTxInputByTxIdAndIdx(ctx *gin.Context) {
	log.Printf("GetTxInputByTxIdAndIdx enter")

	// check tx
	txIdHex := ctx.Param("txid")
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check index
	txIndexString := ctx.Param("index")
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		log.Printf("txindex invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	result, err := service.GetTxInputByTxIdAndIdx(hex.EncodeToString(txId), txIndex)
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txin failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxInputByTxIdAndIdxInsideHeight
// @Summary 通过交易txid和index和交易被打包的区块高度height获取指定交易输入信息
// @Tags Txin
// @Produce  json
// @Param height path int true "Block Height" default(170)
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Param index path int true "input index" default(0)
// @Success 200 {object} model.Response{data=model.TxInResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /height/{height}/tx/{txid}/in/{index} [get]
func GetTxInputByTxIdAndIdxInsideHeight(ctx *gin.Context) {
	log.Printf("GetTxInputByTxIdAndIdxInsideHeight enter")

	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil {
		log.Printf("height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "height invalid"})
		return
	}

	// check tx
	txIdHex := ctx.Param("txid")
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check index
	txIndexString := ctx.Param("index")
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		log.Printf("txindex invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	result, err := service.GetTxInputByTxIdAndIdxInsideHeight(blkHeight, hex.EncodeToString(txId), txIndex)
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txin failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
