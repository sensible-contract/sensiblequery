package controller

import (
	"encoding/hex"
	"net/http"
	"sensiblequery/lib/utils"
	"sensiblequery/logger"
	"sensiblequery/model"
	"sensiblequery/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetTxOutputsByTxId
// @Summary 通过交易txid获取交易所有输出信息列表
// @Tags Txout
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Success 200 {object} model.Response{data=[]model.TxOutStatusResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /tx/{txid}/outs [get]
func GetTxOutputsByTxId(ctx *gin.Context) {
	logger.Log.Info("GetTxOutputsByTxId enter")

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
	if err != nil || size <= 0 {
		logger.Log.Info("size invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
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

	result, err := service.GetTxOutputsByTxId(cursor, size, hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get txouts failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxOutputsByTxIdInsideHeight
// @Summary 通过交易txid和交易被打包的区块高度height获取交易所有输出信息列表
// @Tags Txout
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param height path int true "Block Height" default(170)
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Success 200 {object} model.Response{data=[]model.TxOutStatusResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /height/{height}/tx/{txid}/outs [get]
func GetTxOutputsByTxIdInsideHeight(ctx *gin.Context) {
	logger.Log.Info("GetTxOutputsByTxId enter")

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
	if err != nil || size <= 0 {
		logger.Log.Info("size invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
		return
	}

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

	result, err := service.GetTxOutputsByTxIdInsideHeight(cursor, size, blkHeight, hex.EncodeToString(txId))
	if err != nil {
		logger.Log.Info("get txout with height failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxOutputByTxIdAndIdx
// @Summary 通过交易txid和index获取指定交易输出信息
// @Tags Txout
// @Produce  json
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Param index path int true "output index" default(0)
// @Success 200 {object} model.Response{data=model.TxOutStatusResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Security BearerAuth
// @Router /tx/{txid}/out/{index} [get]
func GetTxOutputByTxIdAndIdx(ctx *gin.Context) {
	logger.Log.Info("GetTxOutputByTxIdAndIdx enter")

	// check tx
	txIdHex := ctx.Param("txid")
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check index
	txIndexString := ctx.Param("index")
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		logger.Log.Info("txindex invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	result, err := service.GetTxOutputByTxIdAndIdx(hex.EncodeToString(txId), txIndex)
	if err != nil {
		logger.Log.Info("get txout failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxOutputByTxIdAndIdxInsideHeight
// @Summary 通过交易txid和index和交易被打包的区块高度height获取指定交易输出信息
// @Tags Txout
// @Produce  json
// @Param height path int true "Block Height" default(170)
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Param index path int true "output index" default(0)
// @Success 200 {object} model.Response{data=model.TxOutStatusResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Security BearerAuth
// @Router /height/{height}/tx/{txid}/out/{index} [get]
func GetTxOutputByTxIdAndIdxInsideHeight(ctx *gin.Context) {
	logger.Log.Info("GetTxOutputByTxIdAndIdxInsideHeight enter")

	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil {
		logger.Log.Info("height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "height invalid"})
		return
	}

	// check tx
	txIdHex := ctx.Param("txid")
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check index
	txIndexString := ctx.Param("index")
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		logger.Log.Info("txindex invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	result, err := service.GetTxOutputByTxIdAndIdxInsideHeight(blkHeight, hex.EncodeToString(txId), txIndex)
	if err != nil {
		logger.Log.Info("get txout with height failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetTxOutputSpentStatusByTxIdAndIdx
// @Summary 通过交易txid和index获取指定交易输出是否被花费状态
// @Tags Txout
// @Produce  json
// @Param txid path string true "TxId" default(f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16)
// @Param index path int true "output index" default(0)
// @Success 200 {object} model.Response{data=model.TxInSpentResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Security BearerAuth
// @Router /tx/{txid}/out/{index}/spent [get]
func GetTxOutputSpentStatusByTxIdAndIdx(ctx *gin.Context) {
	logger.Log.Info("GetTxOutputSpentStatusByTxIdAndIdx enter")

	// check tx
	txIdHex := ctx.Param("txid")
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		logger.Log.Info("txid invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	// check index
	txIndexString := ctx.Param("index")
	txIndex, err := strconv.Atoi(txIndexString)
	if err != nil || txIndex < 0 {
		logger.Log.Info("txindex invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txindex invalid"})
		return
	}

	result, err := service.GetTxOutputSpentStatusByTxIdAndIdx(hex.EncodeToString(txId), txIndex)
	if err != nil {
		logger.Log.Info("get txout spent status failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get vout failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
