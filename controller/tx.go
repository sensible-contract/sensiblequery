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

// GetBlockTxsByBlockHeight
// @Summary 通过区块height获取区块包含的Tx概述列表
// @Tags Tx
// @Produce  json
// @Param height path int true "Block Height" default(3)
// @Success 200 {object} model.Response{data=[]model.TxInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /height/{height}/block/txs [get]
func GetBlockTxsByBlockHeight(ctx *gin.Context) {
	log.Printf("GetBlockTxsByBlockHeight enter")

	// check height
	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil || blkHeight < 0 {
		log.Printf("blk height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk height invalid"})
		return
	}

	blkTxs, err := service.GetBlockTxsByBlockHeight(blkHeight)
	if err != nil {
		log.Printf("get block txs failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get block txs failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: blkTxs,
	})
}

// GetBlockTxsByBlockId
// @Summary 通过区块blkid获取区块包含的Tx概述列表
// @Tags Tx
// @Produce  json
// @Param blkid path string true "Block ID" default(0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449)
// @Success 200 {object} model.Response{data=[]model.TxInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /block/txs/{blkid} [get]
func GetBlockTxsByBlockId(ctx *gin.Context) {
	log.Printf("GetBlockTxsByBlockId enter")

	blkIdHex := ctx.Param("blkid")
	// check
	blkIdReverse, err := hex.DecodeString(blkIdHex)
	if err != nil {
		log.Printf("blkid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blkid invalid"})
		return
	}
	blkId := utils.ReverseBytes(blkIdReverse)

	blkTxs, err := service.GetBlockTxsByBlockId(hex.EncodeToString(blkId))
	if err != nil {
		log.Printf("get block txs failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get block txs failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: blkTxs,
	})
}

// GetTxById
// @Summary 通过交易txid获取交易概述
// @Tags Tx
// @Produce  json
// @Param txid path string true "TxId" default(999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644)
// @Success 200 {object} model.Response{data=model.TxInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /tx/{txid} [get]
func GetTxById(ctx *gin.Context) {
	log.Printf("GetTxById enter")

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetTxById(hex.EncodeToString(txId))
	if err != nil {
		log.Printf("get tx failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: tx,
	})
}

// GetTxByIdInsideHeight
// @Summary 通过交易txid和交易被打包的区块高度height获取交易概述
// @Tags Tx
// @Produce  json
// @Param height path int true "Block Height" default(3)
// @Param txid path string true "TxId" default(999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644)
// @Success 200 {object} model.Response{data=model.TxInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /height/{height}/tx/{txid} [get]
func GetTxByIdInsideHeight(ctx *gin.Context) {
	log.Printf("GetTxByIdInsideHeight enter")

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

	tx, err := service.GetTxByIdInsideHeight(blkHeight, hex.EncodeToString(txId))
	if err != nil {
		log.Printf("get tx failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: tx,
	})
}
