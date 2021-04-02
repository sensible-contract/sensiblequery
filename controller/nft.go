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

// ListAllNFTCodeHash
// @Summary 查询所有NFT CodeHash简述
// @Tags token NFT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.TokenCodeHashResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/codehash/all [get]
func ListAllNFTCodeHash(ctx *gin.Context) {
	log.Printf("ListAllNFTCodeHash enter")

	result, err := service.GetTokenCodeHash(0)
	if err != nil {
		log.Printf("get dummy failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get dummy failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListAllNFTInfo
// @Summary 查询所有NFT Token简述
// @Tags token NFT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.NFTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/info/all [get]
func ListAllNFTInfo(ctx *gin.Context) {
	log.Printf("ListNFTInfo enter")

	result, err := service.GetNFTInfo()
	if err != nil {
		log.Printf("get dummy failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get dummy failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListNFTSummary
// @Summary 查询使用某codehash的NFT Token简述
// @Tags token NFT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Success 200 {object} model.Response{data=[]model.NFTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/codehash-info/{codehash} [get]
func ListNFTSummary(ctx *gin.Context) {
	log.Printf("ListNFTSummary enter")

	codeHashHex := ctx.Param("codehash")
	// check
	_, err := hex.DecodeString(codeHashHex)
	if err != nil {
		log.Printf("codeHash invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	result, err := service.GetNFTSummary(codeHashHex)
	if err != nil {
		log.Printf("get dummy failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get dummy failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetNFTTransferTimesInBlockRange
// @Summary 查询NFT Token在区块中的转移次数，以合约CodeHash+GenesisID，和tokenId来确认一种NFT。
// @Tags token NFT
// @Produce  json
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "Start Block Height" default(3)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param tokenid path int true "Token ID " default(3)
// @Success 200 {object} model.Response{data=[]model.BlockTokenVolumeResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/transfer-volume/{codehash}/{genesis}/{tokenid} [get]
func GetNFTTransferTimesInBlockRange(ctx *gin.Context) {
	log.Printf("GetNFTTransferTimesInBlockRange enter")

	// check height
	blkStartHeightString := ctx.DefaultQuery("start", "0")
	blkStartHeight, err := strconv.Atoi(blkStartHeightString)
	if err != nil || blkStartHeight < 0 {
		log.Printf("blk start height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk start height invalid"})
		return
	}

	blkEndHeightString := ctx.DefaultQuery("end", "0")
	blkEndHeight, err := strconv.Atoi(blkEndHeightString)
	if err != nil || blkEndHeight < 0 {
		log.Printf("blk end height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	if blkEndHeight <= blkStartHeight || (blkEndHeight-blkStartHeight > 1000) {
		log.Printf("blk end height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk end height invalid"})
		return
	}

	codeHashHex := ctx.Param("codehash")
	// check
	_, err = hex.DecodeString(codeHashHex)
	if err != nil {
		log.Printf("codeHash invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "codeHash invalid"})
		return
	}

	genesisIdHex := ctx.Param("genesis")
	// check
	_, err = hex.DecodeString(genesisIdHex)
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

	result, err := service.GetTokenVolumesInBlocksByHeightRange(blkStartHeight, blkEndHeight, codeHashHex, genesisIdHex, 0, tokenIdx)
	if err != nil {
		log.Printf("GetNFTTransferTimesInBlockRange failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "data failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListNFTOwners
// @Summary 查询NFT Token的持有人。获得每个tokenId所属的地址
// @Tags token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.NFTOwnerResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/owners/{codehash}/{genesis} [get]
func ListNFTOwners(ctx *gin.Context) {
	log.Printf("ListNFTOwners enter")

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

	result, err := service.GetNFTOwnersByCodeHashGenesis(cursor, size, codeHash, genesisId)
	if err != nil {
		log.Printf("ListNFTOwners failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "ListNFTOwners failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
	return
}

// ListAllNFTByOwner
// @Summary 查询某人持有的所有NFT Token列表。获得持有的nft数量计数
// @Tags token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.NFTOwnerByAddressResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/summary/{address} [get]
func ListAllNFTByOwner(ctx *gin.Context) {
	log.Printf("ListAllNFTOwners enter")

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

	result, err := service.GetAllNFTBalanceByAddress(cursor, size, addressPkh)
	if err != nil {
		log.Printf("ListAllNFTByOwner failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "ListAllNFTByOwner failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// ListNFTBalanceByOwner
// @Summary 查询某人持有的某NFT Token的所有TokenId
// @Tags token NFT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.NFTOwnerResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/detail/{codehash}/{genesis}/{address} [get]
func ListNFTBalanceByOwner(ctx *gin.Context) {
	log.Printf("ListNFTOwners enter")

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

	result, err := service.GetNFTBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh)
	if err != nil {
		log.Printf("ListNFTBalanceByOwner failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "ListNFTBalanceByOwner failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
