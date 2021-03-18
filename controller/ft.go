package controller

import (
	"log"
	"net/http"
	"satoblock/model"
	"satoblock/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListAllFTInfo
// @Summary 查询所有FT Token简述
// @Tags token FT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.FTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/info/all [get]
func ListAllFTInfo(ctx *gin.Context) {
	log.Printf("ListFTInfo enter")

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

// GetFTTransferVolumeInBlockRange
// @Summary 查询FT Token在区块中的转移数量，以合约CodeHash+GenesisID来确认一种FT
// @Tags token FT
// @Produce  json
// @Param start path int true "Start Block Height" default(0)
// @Param end path int true "Start Block Height" default(3)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.FTTransferVolumeResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/transfer-volume [get]
func GetFTTransferVolumeInBlockRange(ctx *gin.Context) {
	log.Printf("GetFTTransferVolumeInBlockRange enter")

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

// ListFTOwners
// @Summary 查询FT Token的持有人。获得每个地址的token余额
// @Tags token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.FTOwnerBalanceResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/owners/{codehash}/{genesis} [get]
func ListFTOwners(ctx *gin.Context) {
	log.Printf("ListFTOwners enter")

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

// ListAllFTBalanceByOwner
// @Summary 查询某人持有的FT Token列表。获得每个token的余额
// @Tags token FT
// @Produce  json
// @Param cursor query int true "起始游标" default(0)
// @Param size query int true "返回记录数量" default(10)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=[]model.FTOwnerByAddressResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/balance/all/{address} [get]
func ListAllFTBalanceByOwner(ctx *gin.Context) {
	log.Printf("ListAllFTOwners enter")

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

// ListFTBalanceByOwner
// @Summary 查询某人持有的某FT Token的余额
// @Tags token FT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.FTOwnerBalanceResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/balance/{codehash}/{genesis}/{address} [get]
func ListFTBalanceByOwner(ctx *gin.Context) {
	log.Printf("ListFTOwners enter")

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
