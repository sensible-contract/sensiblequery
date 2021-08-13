package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"satosensible/controller"
	_ "satosensible/docs"
	"satosensible/logger"
	"syscall"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

var (
	// 0.0.0.0:8000
	listen_address = os.Getenv("LISTEN")
	is_testnet     = os.Getenv("TESTNET")
)

func KeepJsonContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	}
}

// @title Sensible Browser
// @version 1.0
// @description Sensible 区块浏览器

// @contact.name satosensible
// @contact.url https://github.com/sensing-contract/satosensible
// @contact.email jiedohh@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
func main() {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.Log, true))
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	// go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
	url := ginSwagger.URL("/swagger/doc.json")
	if is_testnet != "" {
		url = ginSwagger.URL("/test/swagger/doc.json")
	}

	store := persistence.NewInMemoryStore(time.Second)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/", controller.Satotx)

	router.POST("/pushtx", controller.PushTx)
	router.POST("/pushtxs", controller.PushTxs)

	router.GET("/blockchain/info", controller.GetBlockchainInfo)

	router.GET("/mempool/info", controller.GetMempoolInfo)

	router.GET("/blocks", controller.GetBlocksByHeightRange)

	router.GET("/block/id/:blkid", controller.GetBlockById)
	router.GET("/block/txs/:blkid", controller.GetBlockTxsByBlockId)

	router.GET("/rawtx/:txid", controller.GetRawTxById)
	router.GET("/relay/:txid", controller.RelayTxById)

	router.GET("/tx/:txid", controller.GetTxById)
	router.GET("/tx/:txid/ins", controller.GetTxInputsByTxId)
	router.GET("/tx/:txid/outs", controller.GetTxOutputsByTxId)

	router.GET("/tx/:txid/in/:index", controller.GetTxInputByTxIdAndIdx)
	router.GET("/tx/:txid/out/:index", controller.GetTxOutputByTxIdAndIdx)

	router.GET("/tx/:txid/out/:index/spent", controller.GetTxOutputSpentStatusByTxIdAndIdx)

	router.GET("/address/:address/utxo",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 1*time.Second, controller.GetUtxoByAddress))

	router.GET("/nft/sell/utxo", controller.GetNFTSellUtxo)
	router.GET("/nft/sell/utxo/:address", controller.GetNFTSellUtxoByAddress)
	router.GET("/nft/sell/utxo/:codehash/:genesis", controller.GetNFTSellUtxoByGenesis)
	router.GET("/nft/sell/utxo/:codehash/:genesis/:token_index", controller.GetNFTSellUtxoDetail)

	router.GET("/ft/utxo/:codehash/:genesis/:address", controller.GetFTUtxo)
	router.GET("/nft/utxo/:codehash/:genesis/:address", controller.GetNFTUtxo)
	router.GET("/nft/utxo-detail/:codehash/:genesis/:token_index", controller.GetNFTUtxoDetailByTokenIndex)

	router.GET("/address/:address/balance", controller.GetBalanceByAddress)

	router.GET("/ft/codehash/all",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.ListAllFTCodeHash))
	router.GET("/ft/codehash-info/:codehash",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.ListFTSummary))
	router.GET("/ft/info/all",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.ListAllFTInfo))
	router.GET("/ft/transfer-times/:codehash/:genesis",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.GetFTTransferVolumeInBlockRange))
	router.GET("/ft/owners/:codehash/:genesis",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 1*time.Second, controller.ListFTOwners))
	router.GET("/ft/summary/:address",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 2*time.Second, controller.ListAllFTBalanceByOwner))
	router.GET("/ft/balance/:codehash/:genesis/:address", controller.GetFTBalanceByOwner)
	router.GET("/ft/history/:codehash/:genesis/:address",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.GetFTHistoryByGenesis))

	router.GET("/nft/codehash/all",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.ListAllNFTCodeHash))
	router.GET("/nft/codehash-info/:codehash",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.ListNFTSummary))
	router.GET("/nft/info/all",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.ListAllNFTInfo))
	router.GET("/nft/transfer-times/:codehash/:genesis/:tokenid",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.GetNFTTransferTimesInBlockRange))
	router.GET("/nft/owners/:codehash/:genesis", controller.ListNFTOwners)
	router.GET("/nft/summary/:address",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 2*time.Second, controller.ListAllNFTByOwner))
	router.GET("/nft/detail/:codehash/:genesis/:address", controller.ListNFTCountByOwner)
	router.GET("/nft/history/:codehash/:genesis/:address",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.GetNFTHistoryByGenesis))

	router.GET("/address/:address/history",
		KeepJsonContentType(), cache.CachePageWithoutHeader(store, 10*time.Second, controller.GetHistoryByAddress))
	router.GET("/contract/history/:codehash/:genesis/:address", controller.GetHistoryByGenesis)

	router.GET("/token/info", controller.ListAllTokenInfo)

	heightAPI := router.Group("/height/:height")
	{
		heightAPI.GET("/block", controller.GetBlockByHeight)

		heightAPI.GET("/block/txs", controller.GetBlockTxsByBlockHeight)

		heightAPI.GET("/rawtx/:txid", controller.GetRawTxByIdInsideHeight)

		heightAPI.GET("/tx/:txid", controller.GetTxByIdInsideHeight)
		heightAPI.GET("/tx/:txid/ins", controller.GetTxInputsByTxIdInsideHeight)
		heightAPI.GET("/tx/:txid/outs", controller.GetTxOutputsByTxIdInsideHeight)

		heightAPI.GET("/tx/:txid/in/:index", controller.GetTxInputByTxIdAndIdxInsideHeight)
		heightAPI.GET("/tx/:txid/out/:index", controller.GetTxOutputByTxIdAndIdxInsideHeight)
	}

	logger.Log.Info("LISTEN:",
		zap.String("address", listen_address),
	)
	svr := &http.Server{
		Addr:    listen_address,
		Handler: router,
	}

	go func() {
		err := svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("ListenAndServe:",
				zap.Error(err),
			)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeout := time.Duration(1) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Shutdown:",
			zap.Error(err),
		)

	}
}
