package controller

import (
	"log"
	"net/http"
	"satoblock/model"
	"satoblock/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListAllNFTInfo
// @Summary 查询所有NFT Token简述
// @Tags token NFT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.NFTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/info/all [get]
func ListAllNFTInfo(ctx *gin.Context) {
	log.Printf("ListNFTInfo enter")

	result, err := service.Dummy()
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
// @Param start path int true "Start Block Height" default(0)
// @Param end path int true "Start Block Height" default(3)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param tokenid path int true "Token ID " default(3)
// @Success 200 {object} model.Response{data=[]model.NFTTransferTimesResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /nft/transfer-times [get]
func GetNFTTransferTimesInBlockRange(ctx *gin.Context) {
	log.Printf("GetNFTTransferTimesInBlockRange enter")

	// check height
	blkStartHeightString := ctx.Param("start")
	blkStartHeight, err := strconv.Atoi(blkStartHeightString)
	if err != nil || blkStartHeight < 0 {
		log.Printf("blk start height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk start height invalid"})
		return
	}
	blkEndHeightString := ctx.Param("end")
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

	result, err := service.Dummy()
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

	result, err := service.Dummy()
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

	result, err := service.Dummy()
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

	result, err := service.Dummy()
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
