package controller

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"satosensible/lib/utils"
	"satosensible/logger"
	"satosensible/model"
	"satosensible/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	is_testnet = os.Getenv("TESTNET")
)

// GetBlockTxsByBlockHeight
// @Summary 通过区块height获取区块包含的Tx概述列表
// @Tags Tx
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param height path int true "Block Height" default(3)
// @Success 200 {object} model.Response{data=[]model.TxInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /height/{height}/block/txs [get]
func GetBlockTxsByBlockHeight(ctx *gin.Context) {
	logger.Log.Info("GetBlockTxsByBlockHeight enter")

	// get cursor/size
	cursorString := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorString)
	if err != nil || cursor < 0 {
		logger.Log.Info("cursor invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "cursor invalid"})
		return
	}
	sizeString := ctx.DefaultQuery("size", "16")
	size, err := strconv.Atoi(sizeString)
	if err != nil || size < 0 {
		logger.Log.Info("size invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
		return
	}

	// check height
	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil || blkHeight < 0 {
		logger.Log.Info("blk height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk height invalid"})
		return
	}

	blkTxs, err := service.GetBlockTxsByBlockHeight(cursor, size, blkHeight)
	if err != nil {
		logger.Log.Info("get block txs failed", zap.Error(err))
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
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param blkid path string true "Block ID" default(0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449)
// @Success 200 {object} model.Response{data=[]model.TxInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /block/txs/{blkid} [get]
func GetBlockTxsByBlockId(ctx *gin.Context) {
	logger.Log.Info("GetBlockTxsByBlockId enter")

	// get cursor/size
	cursorString := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorString)
	if err != nil || cursor < 0 {
		logger.Log.Info("cursor invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "cursor invalid"})
		return
	}
	sizeString := ctx.DefaultQuery("size", "16")
	size, err := strconv.Atoi(sizeString)
	if err != nil || size < 0 {
		logger.Log.Info("size invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
		return
	}

	blkIdHex := ctx.Param("blkid")
	// check
	blkIdReverse, err := hex.DecodeString(blkIdHex)
	if err != nil {
		logger.Log.Info("blkid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blkid invalid"})
		return
	}
	blkId := utils.ReverseBytes(blkIdReverse)

	blkTxs, err := service.GetBlockTxsByBlockId(cursor, size, hex.EncodeToString(blkId))
	if err != nil {
		logger.Log.Info("get block txs failed", zap.Error(err))
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
	logger.Log.Info("GetTxById enter")

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetTxById(hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get tx failed", zap.Error(err))
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
	logger.Log.Info("GetTxByIdInsideHeight enter")

	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil {
		logger.Log.Info("height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "height invalid"})
		return
	}

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetTxByIdInsideHeight(blkHeight, hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get tx failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: tx,
	})
}

////////////////////////////////////////////////////////////////
// GetRawTxById
// @Summary 通过交易txid获取交易原数据rawtx
// @Tags Tx
// @Produce  json
// @Param txid path string true "TxId" default(999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644)
// @Success 200 {object} model.Response{data=string} "{"code": 0, "data": "00...", "msg": "ok"}"
// @Router /rawtx/{txid} [get]
func GetRawTxById(ctx *gin.Context) {
	logger.Log.Info("GetRawTxById enter")

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetRawTxById(hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get tx failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: hex.EncodeToString(tx),
	})
}

////////////////////////////////////////////////////////////////
// RelayTxById
// @Summary 将交易txid重新发送到woc
// @Tags Tx
// @Produce  json
// @Param txid path string true "TxId" default(999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644)
// @Success 200 {object} model.Response{data=string} "{"code": 0, "data": "00...", "msg": "ok"}"
// @Router /relay/{txid} [get]
func RelayTxById(ctx *gin.Context) {
	logger.Log.Info("RelayTxById enter")

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetRawTxById(hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get tx failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	woc := "https://api.whatsonchain.com/v1/bsv/main/tx/raw"
	if is_testnet != "" {
		woc = "https://api.whatsonchain.com/v1/bsv/test/tx/raw"
	}
	jsonData := fmt.Sprintf(`{"txhex": "%s"}`, hex.EncodeToString(tx))
	resp, err := http.Post(woc, "application/json", bytes.NewBufferString(jsonData))
	if err != nil {
		logger.Log.Info("relay tx failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "relay tx failed"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: string(body),
	})
}

// GetRawTxByIdInsideHeight
// @Summary 通过交易txid和交易被打包的区块高度height获取交易原数据rawtx
// @Tags Tx
// @Produce  json
// @Param height path int true "Block Height" default(3)
// @Param txid path string true "TxId" default(999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644)
// @Success 200 {object} model.Response{data=string} "{"code": 0, "data": "00...", "msg": "ok"}"
// @Router /height/{height}/rawtx/{txid} [get]
func GetRawTxByIdInsideHeight(ctx *gin.Context) {
	logger.Log.Info("GetRawTxByIdInsideHeight enter")

	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil {
		logger.Log.Info("height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "height invalid"})
		return
	}

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetRawTxByIdInsideHeight(blkHeight, hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get tx failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: hex.EncodeToString(tx),
	})
}
