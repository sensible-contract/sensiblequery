package controller

import (
	"encoding/hex"
	"net/http"
	"sensiblequery/logger"
	"sensiblequery/model"
	"sensiblequery/service"
	"strconv"

	"github.com/gin-gonic/gin"
	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

// GetContractSwapDataInBlockRange
// @Summary 查询Swap合约在区块中的每次交易数据，以合约CodeHash+GenesisID来确认一种Swap
// @Tags Unique
// @Produce  json
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "End Block Height, (0 to get mempool data)" default(0)
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.ContractSwapDataResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /contract/swap-data/{codehash}/{genesis} [get]
func GetContractSwapDataInBlockRange(ctx *gin.Context) {
	logger.Log.Info("GetContractSwapDataInBlockRange enter")

	// check height
	blkStartHeightString := ctx.DefaultQuery("start", "0")
	blkStartHeight, err := strconv.Atoi(blkStartHeightString)
	if err != nil || blkStartHeight < 0 {
		logger.Log.Info("blk start height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk start height invalid"})
		return
	}
	blkEndHeightString := ctx.DefaultQuery("end", "0")
	blkEndHeight, err := strconv.Atoi(blkEndHeightString)
	if err != nil || blkEndHeight < 0 {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	if blkEndHeight > 0 && (blkEndHeight <= blkStartHeight || (blkEndHeight-blkStartHeight > 10000)) {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

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

	codeHashHex := ctx.Param("codehash")
	// check
	_, err = hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	_, err = hex.DecodeString(genesisIdHex)
	if err != nil {
		logger.Log.Info("genesisId invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetContractSwapDataInBlocksByHeightRange(cursor, size, blkStartHeight, blkEndHeight, codeHashHex, genesisIdHex, scriptDecoder.CodeType_UNIQUE)
	if err != nil {
		logger.Log.Info("get swap failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get swap failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetContractSwapAggregateInBlockRange
// @Summary 查询Swap合约在区块中的聚合价格数据，以合约CodeHash+GenesisID来确认一种Swap
// @Tags Unique
// @Produce  json
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "End Block Height, (0 to get till latest)" default(0)
// @Param interval query int true "每批聚合区块数量" default(6)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.ContractSwapAggregateResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /contract/swap-aggregate/{codehash}/{genesis} [get]
func GetContractSwapAggregateInBlockRange(ctx *gin.Context) {
	logger.Log.Info("GetContractSwapAggregateInBlockRange enter")

	// check height
	blkStartHeightString := ctx.DefaultQuery("start", "0")
	blkStartHeight, err := strconv.Atoi(blkStartHeightString)
	if err != nil || blkStartHeight < 0 {
		logger.Log.Info("blk start height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk start height invalid"})
		return
	}
	blkEndHeightString := ctx.DefaultQuery("end", "0")
	blkEndHeight, err := strconv.Atoi(blkEndHeightString)
	if err != nil || blkEndHeight < 0 {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	if blkEndHeight > 0 && (blkEndHeight <= blkStartHeight || (blkEndHeight-blkStartHeight > 100000)) {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	intervalString := ctx.DefaultQuery("interval", "6")
	interval, err := strconv.Atoi(intervalString)
	if err != nil || interval < 0 {
		logger.Log.Info("interval invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "interval invalid"})
		return
	}

	codeHashHex := ctx.Param("codehash")
	// check
	_, err = hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	_, err = hex.DecodeString(genesisIdHex)
	if err != nil {
		logger.Log.Info("genesisId invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetContractSwapAggregateInBlocksByHeightRange(interval, blkStartHeight, blkEndHeight, codeHashHex, genesisIdHex, scriptDecoder.CodeType_UNIQUE)
	if err != nil {
		logger.Log.Info("get swap failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get swap failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetContractSwapAggregateAmountInBlockRange
// @Summary 查询Swap合约在区块中的聚合Amount数据，以合约CodeHash+GenesisID来确认一种Swap
// @Tags Unique
// @Produce  json
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "End Block Height, (0 to get till latest)" default(0)
// @Param interval query int true "每批聚合区块数量" default(6)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.ContractSwapAggregateAmountResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /contract/swap-aggregate-amount/{codehash}/{genesis} [get]
func GetContractSwapAggregateAmountInBlockRange(ctx *gin.Context) {
	logger.Log.Info("GetContractSwapAggregateAmountInBlockRange enter")

	// check height
	blkStartHeightString := ctx.DefaultQuery("start", "0")
	blkStartHeight, err := strconv.Atoi(blkStartHeightString)
	if err != nil || blkStartHeight < 0 {
		logger.Log.Info("blk start height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk start height invalid"})
		return
	}
	blkEndHeightString := ctx.DefaultQuery("end", "0")
	blkEndHeight, err := strconv.Atoi(blkEndHeightString)
	if err != nil || blkEndHeight < 0 {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	if blkEndHeight > 0 && (blkEndHeight <= blkStartHeight || (blkEndHeight-blkStartHeight > 100000)) {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	intervalString := ctx.DefaultQuery("interval", "6")
	interval, err := strconv.Atoi(intervalString)
	if err != nil || interval < 0 {
		logger.Log.Info("interval invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "interval invalid"})
		return
	}

	codeHashHex := ctx.Param("codehash")
	// check
	_, err = hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	_, err = hex.DecodeString(genesisIdHex)
	if err != nil {
		logger.Log.Info("genesisId invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetContractSwapAggregateAmountInBlocksByHeightRange(interval, blkStartHeight, blkEndHeight, codeHashHex, genesisIdHex, scriptDecoder.CodeType_UNIQUE)
	if err != nil {
		logger.Log.Info("get swap failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get swap failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
