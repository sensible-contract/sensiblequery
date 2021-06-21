package controller

import (
	"encoding/hex"
	"log"
	"net/http"
	"satosensible/lib/utils"
	"satosensible/model"
	"satosensible/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBalanceByAddress
// @Summary 通过地址address获取balance
// @Tags UTXO
// @Produce  json
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.BalanceResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /address/{address}/balance [get]
func GetBalanceByAddress(ctx *gin.Context) {
	log.Printf("GetBalanceByAddress enter")

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		log.Printf("address invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}
	log.Printf("address: %s", hex.EncodeToString(addressPkh))
	result, err := service.GetBalanceByAddress(addressPkh)
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
	log.Printf("GetUtxoByAddress enter")

	// get cursor/size
	cursorString := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorString)
	if err != nil || cursor < 0 {
		log.Printf("cursor invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "cursor invalid"})
		return
	}
	sizeString := ctx.DefaultQuery("size", "16")
	size, err := strconv.Atoi(sizeString)
	if err != nil || size < 0 {
		log.Printf("size invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
		return
	}

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		log.Printf("address invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, err := service.GetUtxoByAddress(cursor, size, addressPkh)
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

// GetNFTUtxoDetailByTokenId
// @Summary 通过NFT合约CodeHash+溯源genesis获取某tokenId的utxo
// @Tags UTXO, token NFT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID" default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param tokenid path int true "Token ID" default(3)
// @Success 200 {object} model.Response{data=[]model.TxOutResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/utxo-detail/{codehash}/{genesis}/{tokenid} [get]
func GetNFTUtxoDetailByTokenId(ctx *gin.Context) {
	log.Printf("GetNFTUtxoDetailByTokenId enter")

	codeHashHex := ctx.Param("codehash")
	// check
	codeHash, err := hex.DecodeString(codeHashHex)
	if err != nil {
		log.Printf("codeHash invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		log.Printf("genesisId invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	tokenIdxString := ctx.Param("tokenid")
	tokenIdx, err := strconv.Atoi(tokenIdxString)
	if err != nil || tokenIdx < 0 {
		log.Printf("tokenIdx invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tokenIdx invalid"})
		return
	}

	result, err := service.GetUtxoByTokenId(codeHash, genesisId, tokenIdxString)
	if err != nil {
		log.Printf("get nft utxo detail failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get txo failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
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
	log.Printf("GetFTUtxo enter")
	GetUtxoByCodeHashGenesisAddress(ctx, false)
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
	log.Printf("GetNFTUtxo enter")
	GetUtxoByCodeHashGenesisAddress(ctx, true)
}

func GetUtxoByCodeHashGenesisAddress(ctx *gin.Context, isNFT bool) {
	// get cursor/size
	cursorString := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorString)
	if err != nil || cursor < 0 {
		log.Printf("cursor invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "cursor invalid"})
		return
	}
	sizeString := ctx.DefaultQuery("size", "16")
	size, err := strconv.Atoi(sizeString)
	if err != nil || size < 0 {
		log.Printf("size invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "size invalid"})
		return
	}

	codeHashHex := ctx.Param("codehash")
	// check
	codeHash, err := hex.DecodeString(codeHashHex)
	if err != nil {
		log.Printf("codeHash invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	genesisId, err := hex.DecodeString(genesisIdHex)
	if err != nil {
		log.Printf("genesisId invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "genesisId invalid"})
		return
	}

	address := ctx.Param("address")
	// check
	addressPkh, err := utils.DecodeAddress(address)
	if err != nil {
		log.Printf("address invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "address invalid"})
		return
	}

	result, err := service.GetUtxoByCodeHashGenesisAddress(cursor, size, codeHash, genesisId, addressPkh, isNFT)
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
