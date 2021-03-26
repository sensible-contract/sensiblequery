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

// ListAllFTInfo
// @Summary 查询所有FT Token简述
// @Tags token FT
// @Produce  json
// @Success 200 {object} model.Response{data=[]model.FTInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/info/all [get]
func ListAllFTInfo(ctx *gin.Context) {
	log.Printf("ListFTInfo enter")

	result, err := service.GetFTInfo()
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
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "Start Block Height" default(3)
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Success 200 {object} model.Response{data=[]model.BlockTokenVolumeResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/transfer-volume/{codehash}/{genesis} [get]
func GetFTTransferVolumeInBlockRange(ctx *gin.Context) {
	log.Printf("GetFTTransferVolumeInBlockRange enter")

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

	result, err := service.GetTokenVolumesInBlocksByHeightRange(blkStartHeight, blkEndHeight, codeHashHex, genesisIdHex, 0, 0)
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

	result, err := service.GetTokenOwnersByCodeHashGenesis(cursor, size, codeHash, genesisId)
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
// @Router /ft/summary/{address} [get]
func ListAllFTBalanceByOwner(ctx *gin.Context) {
	log.Printf("ListAllFTOwners enter")
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

	result, err := service.GetAllTokenBalanceByAddress(cursor, size, addressPkh)
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

// GetFTBalanceByOwner
// @Summary 查询某人持有的某FT Token的余额
// @Tags token FT
// @Produce  json
// @Param codehash path string true "Code Hash160" default(844c56bb99afc374967a27ce3b46244e2e1fba60)
// @Param genesis path string true "Genesis ID " default(74967a27ce3b46244e2e1fba60844c56bb99afc3)
// @Param address path string true "Address" default(17SkEw2md5avVNyYgj6RiXuQKNwkXaxFyQ)
// @Success 200 {object} model.Response{data=model.FTOwnerBalanceResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /ft/balance/{codehash}/{genesis}/{address} [get]
func GetFTBalanceByOwner(ctx *gin.Context) {
	log.Printf("GetFTBalanceByOwner enter")
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

	result, err := service.GetTokenBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh)
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
