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
	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

// ListAllFTCodeHash
// @Summary 查询所有FT codehash简述
// @Tags token FT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.TokenCodeHashResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/codehash/all [get]
func ListAllFTCodeHash(ctx *gin.Context) {
	logger.Log.Info("ListAllFTCodeHash enter")

	result, err := service.GetTokenCodeHash(scriptDecoder.CodeType_FT)
	if err != nil {
		logger.Log.Info("get ft codehash failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get ft codehash failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListAllFTInfo
// @Summary 查询所有FT Token简述
// @Tags token FT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.FTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/info/all [get]
func ListAllFTInfo(ctx *gin.Context) {
	logger.Log.Info("ListFTInfo enter")

	result, err := service.GetFTInfo()
	if err != nil {
		logger.Log.Info("get ft info failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get ft info failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListFTSummary
// @Summary 查询使用某codehash的FT Token简述
// @Tags token FT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Success 200 {object} model.Response{data=[]model.FTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/codehash-info/{codehash} [get]
func ListFTSummary(ctx *gin.Context) {
	logger.Log.Info("ListFTSummary enter")

	codeHashHex := ctx.Param("codehash")
	// check
	_, err := hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	result, err := service.GetFTSummary(codeHashHex)
	if err != nil {
		logger.Log.Info("get ft summary failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get ft summary failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListFTInfoByGenesis
// @Summary 查询使用某codehash+genesis的FT Token简述
// @Tags token FT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=model.FTInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /ft/genesis-info/{codehash}/{genesis} [get]
func ListFTInfoByGenesis(ctx *gin.Context) {
	logger.Log.Info("ListFTInfoByGenesis enter")

	codeHashHex := ctx.Param("codehash")
	// check
	_, err := hex.DecodeString(codeHashHex)
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

	result, err := service.ListFTInfoByGenesis(codeHashHex, genesisIdHex)
	if err != nil {
		logger.Log.Info("get ft summary failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get ft summary failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetFTTransferVolumeInBlockRange
// @Summary 查询FT Token在区块中的转移数量，以合约CodeHash+GenesisID来确认一种FT
// @Tags token FT
// @Produce  json
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "End Block Height" default(0)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.BlockTokenVolumeResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/transfer-times/{codehash}/{genesis} [get]
func GetFTTransferVolumeInBlockRange(ctx *gin.Context) {
	logger.Log.Info("GetFTTransferVolumeInBlockRange enter")

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

	if blkEndHeight <= blkStartHeight || (blkEndHeight-blkStartHeight > 1000) {
		logger.Log.Info("blk end height invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
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

	result, err := service.GetTokenVolumesInBlocksByHeightRange(blkStartHeight, blkEndHeight, codeHashHex, genesisIdHex, scriptDecoder.CodeType_FT, 0)
	if err != nil {
		logger.Log.Info("get token volumes failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get token volumes failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListFTOwners
// @Summary 查询FT Token的持有人。获得每个地址的token余额
// @Tags token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.FTOwnerBalanceResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/owners/{codehash}/{genesis} [get]
func ListFTOwners(ctx *gin.Context) {
	logger.Log.Info("ListFTOwners enter")

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
	codeHash, err := hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		logger.Log.Info("genesisId invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	result, err := service.GetTokenOwnersByCodeHashGenesis(cursor, size, codeHash, genesisId)
	if err != nil {
		logger.Log.Info("get token owner failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get token owner failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListAllFTSummaryDataByOwner
// @Summary 查询某人持有的FT Token列表。获得每个token的余额。并返回token总量
// @Tags token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.FTSummaryDataByAddressResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /ft/summary-data/{address} [get]
func ListAllFTSummaryDataByOwner(ctx *gin.Context) {
	logger.Log.Info("ListAllFTOwners enter")

	ListAllFTSummaryByOwnerCommon(ctx, true)
}

// ListAllFTSummaryByOwner
// @Summary 查询某人持有的FT Token列表。获得每个token的余额
// @Tags token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.FTSummaryByAddressResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/summary/{address} [get]
func ListAllFTSummaryByOwner(ctx *gin.Context) {
	logger.Log.Info("ListAllFTOwners enter")

	ListAllFTSummaryByOwnerCommon(ctx, false)
}

func ListAllFTSummaryByOwnerCommon(ctx *gin.Context, page bool) {
	logger.Log.Info("ListAllFTOwners enter")
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

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		logger.Log.Info("address invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, total, err := service.GetAllTokenBalanceByAddress(cursor, size, addressPkh)
	if err != nil {
		logger.Log.Info("get token balance failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get token balance failed"})
		return
	}

	if page {
		ctx.JSON(http.StatusOK, model.Response{
			Code: 0,
			Msg:  "ok",
			Data: &model.FTSummaryDataByAddressResp{
				Cursor: cursor,
				Total:  total,
				Token:  result,
			},
		})

	} else {
		ctx.JSON(http.StatusOK, model.Response{
			Code: 0,
			Msg:  "ok",
			Data: result,
		})
	}
}

// GetFTBalanceByOwner
// @Summary 查询某人持有的某FT Token的余额，同时返回UTXO数量
// @Tags token FT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.FTOwnerBalanceWithUtxoCountResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/balance/{codehash}/{genesis}/{address} [get]
func GetFTBalanceByOwner(ctx *gin.Context) {
	logger.Log.Info("GetFTBalanceByOwner enter")
	codeHashHex := ctx.Param("codehash")
	// check
	codeHash, err := hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		logger.Log.Info("genesisId invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		logger.Log.Info("address invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, err := service.GetTokenBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh)
	if err != nil {
		logger.Log.Info("get ft balance failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get ft balance failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
