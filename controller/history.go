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

const MAX_HISTORY_LIMIT = 1024000

// GetHistoryByAddress
// @Summary 通过地址address获取相关tx历史列表，返回详细输入/输出
// @Tags History
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /address/{address}/history [get]
func GetHistoryByAddress(ctx *gin.Context) {
	logger.Log.Info("GetHistoryByAddress enter")
	GetHistoryByAddressAndType(ctx, model.HISTORY_CONTRACT_P2PKH_BOTH)
}

// GetContractHistoryByAddress
// @Summary 通过地址address获取合约相关tx历史列表，返回详细输入/输出
// @Tags History
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /address/{address}/contract-history [get]
func GetContractHistoryByAddress(ctx *gin.Context) {
	logger.Log.Info("GetContractHistoryByAddress enter")
	GetHistoryByAddressAndType(ctx, model.HISTORY_CONTRACT_ONLY)
}

func GetHistoryByAddressAndType(ctx *gin.Context, historyType model.HistoryType) {
	logger.Log.Info("GetHistoryByAddressAndType enter")
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
	if err != nil || size <= 0 || cursor+size > MAX_HISTORY_LIMIT {
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

	result, err := service.GetHistoryByAddressAndType(cursor, size, hex.EncodeToString(addressPkh), historyType)
	if err != nil {
		logger.Log.Info("get history failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get history failed"})
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
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutHistoryResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /contract/history/{codehash}/{genesis}/{address} [get]
func GetHistoryByGenesis(ctx *gin.Context) {
	logger.Log.Info("GetHistoryByGenesis enter")

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
	if err != nil || size <= 0 || cursor+size > MAX_HISTORY_LIMIT {
		logger.Log.Info("size invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
		return
	}

	codehashHex := ctx.Param("codehash")
	// check
	_, err = hex.DecodeString(codehashHex)
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

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		logger.Log.Info("address invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}
	result, err := service.GetHistoryByGenesis(cursor, size, codehashHex, genesisIdHex, hex.EncodeToString(addressPkh))
	if err != nil {
		logger.Log.Info("get history failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get histroy failed"})
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
	logger.Log.Info("GetFTHistoryByGenesis enter")
	GetHistoryByGenesis(ctx)
}

// GetNFTHistoryByGenesis
// @Summary 通过NFT合约CodeHash+溯源genesis获取地址相关tx历史列表，返回详细输入/输出
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
	logger.Log.Info("GetNFTHistoryByGenesis enter")
	GetHistoryByGenesis(ctx)
}
