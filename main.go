package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"satoblock/controller"
	_ "satoblock/docs"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	// 0.0.0.0:8000
	listen_address = os.Getenv("LISTEN")
)

// @title Sensible Browser
// @version 1.0
// @description Sensible 区块浏览器

// @contact.name satoblock
// @contact.url https://github.com/sensing-contract/satoblock
// @contact.email jiedohh@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/", controller.Satotx)

	router.POST("/pushtx", controller.PushTx)

	router.GET("/blockchain/info", controller.GetBlockchainInfo)

	router.GET("/blocks", controller.GetBlocksByHeightRange)

	router.GET("/block/id/:blkid", controller.GetBlockById)
	router.GET("/block/txs/:blkid", controller.GetBlockTxsByBlockId)

	router.GET("/tx/:txid", controller.GetTxById)
	router.GET("/tx/:txid/ins", controller.GetTxInputsByTxId)
	router.GET("/tx/:txid/outs", controller.GetTxOutputsByTxId)

	router.GET("/tx/:txid/in/:index", controller.GetTxInputByTxIdAndIdx)
	router.GET("/tx/:txid/out/:index", controller.GetTxOutputByTxIdAndIdx)

	router.GET("/tx/:txid/out/:index/spent", controller.GetTxOutputSpentStatusByTxIdAndIdx)

	router.GET("/address/:address/utxo", controller.GetUtxoByAddress)

	router.GET("/ft/utxo/:codehash/:genesis/:address", controller.GetFTUtxo)
	router.GET("/nft/utxo/:codehash/:genesis/:address", controller.GetNFTUtxo)

	router.GET("/address/:address/balance", controller.GetBalanceByAddress)

	router.GET("/ft/codehash/all", controller.ListAllFTCodeHash)
	router.GET("/ft/codehash-info/:codehash", controller.ListFTSummary)
	router.GET("/ft/info/all", controller.ListAllFTInfo)
	router.GET("/ft/transfer-volume/:codehash/:genesis", controller.GetFTTransferVolumeInBlockRange)
	router.GET("/ft/owners/:codehash/:genesis", controller.ListFTOwners)
	router.GET("/ft/summary/:address", controller.ListAllFTBalanceByOwner)
	router.GET("/ft/balance/:codehash/:genesis/:address", controller.GetFTBalanceByOwner)

	router.GET("/nft/codehash/all", controller.ListAllNFTCodeHash)
	router.GET("/nft/codehash-info/:codehash", controller.ListNFTSummary)
	router.GET("/nft/info/all", controller.ListAllNFTInfo)
	router.GET("/nft/transfer-volume/:codehash/:genesis/:tokenid", controller.GetNFTTransferTimesInBlockRange)
	router.GET("/nft/owners/:codehash/:genesis", controller.ListNFTOwners)
	router.GET("/nft/summary/:address", controller.ListAllNFTByOwner)
	router.GET("/nft/detail/:codehash/:genesis/:address", controller.ListNFTCountByOwner)

	router.GET("/address/:address/history", controller.GetHistoryByAddress)
	router.GET("/contract/history/:codehash/:genesis/:address", controller.GetHistoryByGenesis)

	router.GET("/token/info", controller.ListAllTokenInfo)

	heightAPI := router.Group("/height/:height")
	{
		heightAPI.GET("/block", controller.GetBlockByHeight)

		heightAPI.GET("/block/txs", controller.GetBlockTxsByBlockHeight)

		heightAPI.GET("/tx/:txid", controller.GetTxByIdInsideHeight)
		heightAPI.GET("/tx/:txid/ins", controller.GetTxInputsByTxIdInsideHeight)
		heightAPI.GET("/tx/:txid/outs", controller.GetTxOutputsByTxIdInsideHeight)

		heightAPI.GET("/tx/:txid/in/:index", controller.GetTxInputByTxIdAndIdxInsideHeight)
		heightAPI.GET("/tx/:txid/out/:index", controller.GetTxOutputByTxIdAndIdxInsideHeight)
	}

	log.Printf("LISTEN: %s", listen_address)
	svr := &http.Server{
		Addr:    listen_address,
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(1) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
