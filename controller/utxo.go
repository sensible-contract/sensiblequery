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

const MAX_UTXO_LIMIT = 5120

// GetBalanceByAddress
// @Summary 通过地址address获取balance
// @Tags UTXO
// @Produce  json
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.BalanceResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /address/{address}/balance [get]
func GetBalanceByAddress(ctx *gin.Context) {
	logger.Log.Info("GetBalanceByAddress enter")

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		logger.Log.Info("address invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}
	logger.Log.Info("GetBalance", zap.String("address", hex.EncodeToString(addressPkh)))
	result, err := service.GetBalanceByAddress(addressPkh)
	if err != nil {
		logger.Log.Info("get balance failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetUtxoDataByAddress
// @Summary 通过地址address获取相关常规utxo列表，和数量信息
// @Tags UTXO
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.AddressUTXOResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /address/{address}/utxo-data [get]
func GetUtxoDataByAddress(ctx *gin.Context) {
	logger.Log.Info("GetUtxoDataByAddress enter")
	GetUtxoDataByAddressCommon(ctx, true)
}

// GetUtxoByAddress
// @Summary 通过地址address获取相关常规utxo列表
// @Tags UTXO
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(16)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxStandardOutResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /address/{address}/utxo [get]
func GetUtxoByAddress(ctx *gin.Context) {
	logger.Log.Info("GetUtxoByAddress enter")
	GetUtxoDataByAddressCommon(ctx, false)
}

func GetUtxoDataByAddressCommon(ctx *gin.Context, detail bool) {
	logger.Log.Info("GetUtxoDataByAddressCommon enter")
	// get cursor
	cursorString := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorString)
	if err != nil || cursor < 0 {
		logger.Log.Info("cursor invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "cursor invalid"})
		return
	}

	// get size
	sizeString := ctx.DefaultQuery("size", "16")
	size, err := strconv.Atoi(sizeString)
	if err != nil || size <= 0 || cursor+size > MAX_UTXO_LIMIT {
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

	result, total, totalConf, totalUnconf, totalUnconfSpend, err := service.GetUtxoByAddress(cursor, size, addressPkh)
	if err != nil {
		logger.Log.Info("get utxo failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	if detail {
		ctx.JSON(http.StatusOK, model.Response{
			Code: 0,
			Msg:  "ok",
			Data: &model.AddressUTXOResp{
				Cursor:                cursor,
				Total:                 total,
				TotalConfirmed:        totalConf,
				TotalUnconfirmedNew:   totalUnconf,
				TotalUnconfirmedSpend: totalUnconfSpend,
				UTXO:                  result,
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

// GetNFTUtxoDetailByTokenIndex
// @Summary 通过NFT合约CodeHash+溯源genesis获取某tokenId的utxo
// @Tags UTXO, token NFT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param token_index path int true "Token Index" default(3)
// @Success 200 {object} model.Response{data=model.TxOutResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /nft/utxo-detail/{codehash}/{genesis}/{token_index} [get]
func GetNFTUtxoDetailByTokenIndex(ctx *gin.Context) {
	logger.Log.Info("GetNFTUtxoDetailByTokenIndex enter")

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

	tokenIndexString := ctx.Param("token_index")
	tokenIndex, err := strconv.Atoi(tokenIndexString)
	if err != nil || tokenIndex < 0 {
		logger.Log.Info("tokenIndex invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tokenIndex invalid"})
		return
	}

	result, err := service.GetUtxoByTokenIndex(codeHash, genesisId, tokenIndexString)
	if err != nil {
		logger.Log.Info("get nft utxo detail failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetFTUtxoData
// @Summary 通过FT合约CodeHash+溯源genesis获取某地址的utxo列表，和数量信息
// @Tags UTXO, token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.AddressTokenUTXOResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /ft/utxo-data/{codehash}/{genesis}/{address} [get]
func GetFTUtxoData(ctx *gin.Context) {
	logger.Log.Info("GetFTUtxoData enter")
	GetUtxoByCodeHashGenesisAddress(ctx, "fu", true)
}

// GetNFTUtxoData
// @Summary 通过NFT合约CodeHash+溯源genesis获取某地址的utxo列表，和数量信息
// @Tags UTXO, token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.AddressTokenUTXOResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /nft/utxo-data/{codehash}/{genesis}/{address} [get]
func GetNFTUtxoData(ctx *gin.Context) {
	logger.Log.Info("GetNFTUtxoData enter")
	GetUtxoByCodeHashGenesisAddress(ctx, "nu", true)
}

// GetFTUtxo
// @Summary 通过FT合约CodeHash+溯源genesis获取某地址的utxo列表
// @Tags UTXO, token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/utxo/{codehash}/{genesis}/{address} [get]
func GetFTUtxo(ctx *gin.Context) {
	logger.Log.Info("GetFTUtxo enter")
	GetUtxoByCodeHashGenesisAddress(ctx, "fu", false)
}

// GetNFTUtxo
// @Summary 通过NFT合约CodeHash+溯源genesis获取某地址的utxo列表
// @Tags UTXO, token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.TxOutResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/utxo/{codehash}/{genesis}/{address} [get]
func GetNFTUtxo(ctx *gin.Context) {
	logger.Log.Info("GetNFTUtxo enter")
	GetUtxoByCodeHashGenesisAddress(ctx, "nu", false)
}

func GetUtxoByCodeHashGenesisAddress(ctx *gin.Context, key string, detail bool) {
	logger.Log.Info("GetUtxoByCodeHashGenesisAddress enter")

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

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		logger.Log.Info("address invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, total, totalConf, totalUnconf, totalUnconfSpend, err := service.GetUtxoByCodeHashGenesisAddress(cursor, size, codeHash, genesisId, addressPkh, key)
	if err != nil {
		logger.Log.Info("get token utxo failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	if detail {
		ctx.JSON(http.StatusOK, model.Response{
			Code: 0,
			Msg:  "ok",
			Data: &model.AddressTokenUTXOResp{
				Cursor:                cursor,
				Total:                 total,
				TotalConfirmed:        totalConf,
				TotalUnconfirmedNew:   totalUnconf,
				TotalUnconfirmedSpend: totalUnconfSpend,
				UTXO:                  result,
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
