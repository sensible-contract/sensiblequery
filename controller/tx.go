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

func GetBlockTxsByBlockHeight(ctx *gin.Context) {
	log.Printf("GetBlockTxsByBlockHeight enter")

	// check height
	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil || blkHeight < 0 {
		log.Printf("blk height invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blk height invalid"})
		return
	}

	blkTxs, err := service.GetBlockTxsByBlockHeight(blkHeight)
	if err != nil {
		log.Printf("get block txs failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get block txs failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: blkTxs,
	})
}

func GetBlockTxsByBlockId(ctx *gin.Context) {
	log.Printf("GetBlockTxsByBlockId enter")

	blkIdHex := ctx.Param("blkid")
	// check
	blkIdReverse, err := hex.DecodeString(blkIdHex)
	if err != nil {
		log.Printf("blkid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "blkid invalid"})
		return
	}
	blkId := utils.ReverseBytes(blkIdReverse)

	blkTxs, err := service.GetBlockTxsByBlockId(hex.EncodeToString(blkId))
	if err != nil {
		log.Printf("get block txs failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get block txs failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: blkTxs,
	})
}

func GetTxById(ctx *gin.Context) {
	log.Printf("GetTxById enter")

	// may have height
	blkHeightString := ctx.Param("height")
	blkHeight, err := strconv.Atoi(blkHeightString)
	if err != nil {
		blkHeight = -1
	}

	txIdHex := ctx.Param("txid")
	// check
	txIdReverse, err := hex.DecodeString(txIdHex)
	if err != nil {
		log.Printf("txid invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "txid invalid"})
		return
	}
	txId := utils.ReverseBytes(txIdReverse)

	tx, err := service.GetTxById(blkHeight, hex.EncodeToString(txId))
	if err != nil {
		log.Printf("get tx failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "get tx failed"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: tx,
	})
}
