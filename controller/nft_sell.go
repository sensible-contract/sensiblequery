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

// GetNFTSellUtxo
// @Summary 获取NFTSell合约的utxo列表
// @Tags UTXO, token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Success 200 {object} model.Response{data=[]model.NFTSellResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /nft/sell/utxo [get]
func GetNFTSellUtxo(ctx *gin.Context) {
	logger.Log.Info("GetNFTSellUtxo enter")

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

	result, err := service.GetNFTSellUtxo(cursor, size)
	if err != nil {
		logger.Log.Info("get block failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetNFTSellUtxoByAddress
// @Summary 通过出售人地址获取NFTSell合约utxo列表
// @Tags UTXO, token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.NFTSellResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /nft/sell/utxo-by-address/{address} [get]
func GetNFTSellUtxoByAddress(ctx *gin.Context) {
	logger.Log.Info("GetNFTSellUtxoByAddress enter")

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

	result, err := service.GetNFTSellUtxoByAddress(cursor, size, addressPkh)
	if err != nil {
		logger.Log.Info("get block failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})

}

// GetNFTSellUtxoByGenesis
// @Summary 通过NFT的CodeHash+溯源genesis获取NFTSell合约utxo列表
// @Tags UTXO, token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.NFTSellResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /nft/sell/utxo/{codehash}/{genesis} [get]
func GetNFTSellUtxoByGenesis(ctx *gin.Context) {
	logger.Log.Info("GetNFTSellUtxoByGenesis enter")

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

	result, err := service.GetNFTSellUtxoByGenesis(cursor, size, codeHash, genesisId)
	if err != nil {
		logger.Log.Info("get block failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})

}

// GetNFTSellUtxoDetail
// @Summary 通过NFT的CodeHash+溯源genesis和NFT Token Index获取具体NFTSell合约utxo
// @Tags UTXO, token NFT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param token_index path int true "Token Index" default(3)
// @Param ready query boolean true "仅返回ready状态的记录" default(true)
// @Success 200 {object} model.Response{data=[]model.NFTSellResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /nft/sell/utxo-detail/{codehash}/{genesis}/{token_index} [get]
func GetNFTSellUtxoDetail(ctx *gin.Context) {
	logger.Log.Info("GetNFTSellUtxoDetail enter")

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

	isReadyOnly := (ctx.DefaultQuery("ready", "true") == "true")
	result, err := service.GetNFTSellUtxoByTokenIndexMerge(codeHash, genesisId, tokenIndexString, isReadyOnly)
	if err != nil {
		logger.Log.Info("GetNFTSellUtxoByTokenIndexMerge", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetNFTAuctionUtxoDetail
// @Summary 通过拍卖的CodeHash和NFT ID获取具体NFTAuction合约utxo
// @Tags UTXO, token NFT
// @Produce  json
// @Param codehash path string true "Auction Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param nftid path string true "NFT ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param ready query boolean true "仅返回ready状态的记录" default(true)
// @Success 200 {object} model.Response{data=[]model.NFTAuctionResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Security BearerAuth
// @Router /nft/auction/utxo-detail/{codehash}/{nftid} [get]
func GetNFTAuctionUtxoDetail(ctx *gin.Context) {
	logger.Log.Info("GetNFTAuctionUtxoDetail enter")

	codeHashHex := ctx.Param("codehash")
	// check
	codeHash, err := hex.DecodeString(codeHashHex)
	if err != nil {
		logger.Log.Info("codeHash invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	nftIdHex := ctx.Param("nftid")
	// check
	nftId, err := hex.DecodeString(nftIdHex)
	if err != nil {
		logger.Log.Info("nftId invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "nftId invalid"})
		return
	}

	isReadyOnly := (ctx.DefaultQuery("ready", "true") == "true")
	result, err := service.GetNFTAuctionUtxoByNFTIDMerge(codeHash, nftId, isReadyOnly)
	if err != nil {
		logger.Log.Info("GetNFTAuctionUtxoByNFTIDMerge", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
