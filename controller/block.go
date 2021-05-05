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

// GetBlocksByHeightRange
// @Summary 获取指定高度范围内的区块概述列表
// @Tags Block
// @Produce  json
// @Param start query int true "Start Block Height" default(0)
// @Param end query int true "Start Block Height" default(3)
// @Success 200 {object} model.Response{data=[]model.BlockInfoResp} "{"code": 0, "data": [{}], "msg": "ok"}"
// @Router /blocks [get]
func GetBlocksByHeightRange(ctx *gin.Context) {
	log.Printf("GetBlocksByHeightRange enter")

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

	result, err := service.GetBlocksByHeightRange(blkStartHeight, blkEndHeight)
	if err != nil {
		log.Printf("get blocks failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get blocks failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetBlockByHeight
// @Summary 通过区块height获取区块概述
// @Tags Block
// @Produce  json
// @Param height path int true "Block Height" default(0)
// @Success 200 {object} model.Response{data=model.BlockInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /height/{height}/block [get]
func GetBlockByHeight(ctx *gin.Context) {
	log.Printf("GetBlockByHeight enter")

	// check height
	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil || blkHeight < 0 {
		log.Printf("blk height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk height invalid"})
		return
	}

	result, err := service.GetBlockByHeight(blkHeight)
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get block failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// GetBlockById
// @Summary 通过区块blkid获取区块概述
// @Tags Block
// @Produce  json
// @Param blkid path string true "BlockId" default(0000000082b5015589a3fdf2d4baff403e6f0be035a5d9742c1cae6295464449)
// @Success 200 {object} model.Response{data=model.BlockInfoResp} "{"code": 0, "data": {}, "msg": "ok"}"
// @Router /block/id/{blkid} [get]
func GetBlockById(ctx *gin.Context) {
	log.Printf("GetBlockById enter")

	blkIdHex := ctx.Param("blkid")
	// check
	blkIdReverse, err := hex.DecodeString(blkIdHex)
	if err != nil {
		log.Printf("blkid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blkid invalid"})
		return
	}
	blkId := utils.ReverseBytes(blkIdReverse)

	result, err := service.GetBlockById(hex.EncodeToString(blkId))
	if err != nil {
		log.Printf("get block failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get block failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}
