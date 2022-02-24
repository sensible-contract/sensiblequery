package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sensiblequery/controller"
	"sensiblequery/dao/rdb"
	_ "sensiblequery/docs"
	"sensiblequery/logger"
	"syscall"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"

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
	basePath       = os.Getenv("BASE_PATH")
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

// @contact.name sensiblequery
// @contact.url https://github.com/sensing-contract/sensiblequery
// @contact.email jiedohh@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
func main() {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.Log, true))
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	// go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
	url := ginSwagger.URL(basePath + "/swagger/doc.json")

	// store := persist.NewMemoryStore(time.Second)
	store := persist.NewRedisStore(rdb.CacheClient)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/", controller.Satotx)

	router.GET("/getrawmempool", controller.GetRawMempool)
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
		cache.CacheByRequestURI(store, 1*time.Second), controller.GetUtxoByAddress)
	router.GET("/address/:address/utxo-data",
		cache.CacheByRequestURI(store, 1*time.Second), controller.GetUtxoDataByAddress)

	router.GET("/nft/sell/utxo", controller.GetNFTSellUtxo)
	router.GET("/nft/sell/utxo-by-address/:address", controller.GetNFTSellUtxoByAddress)
	router.GET("/nft/sell/utxo/:codehash/:genesis", controller.GetNFTSellUtxoByGenesis)
	router.GET("/nft/sell/utxo-detail/:codehash/:genesis/:token_index", controller.GetNFTSellUtxoDetail)

	router.GET("/nft/auction/utxo-detail/:codehash/:nftid", controller.GetNFTAuctionUtxoDetail)

	router.GET("/ft/utxo-data/:codehash/:genesis/:address", controller.GetFTUtxoData)
	router.GET("/nft/utxo-data/:codehash/:genesis/:address", controller.GetNFTUtxoData)

	router.GET("/ft/utxo/:codehash/:genesis/:address", controller.GetFTUtxo)
	router.GET("/nft/utxo/:codehash/:genesis/:address", controller.GetNFTUtxo)
	router.GET("/nft/utxo-detail/:codehash/:genesis/:token_index", controller.GetNFTUtxoDetailByTokenIndex)
	router.GET("/nft/utxo-list/:codehash/:genesis",
		cache.CacheByRequestURI(store, 1*time.Second), controller.GetNFTUtxoList)

	router.GET("/address/:address/balance", controller.GetBalanceByAddress)

	router.GET("/contract/swap-data/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetContractSwapDataInBlockRange)
	router.GET("/contract/swap-aggregate/:codehash/:genesis",
		cache.CacheByRequestURI(store, 60*time.Second), controller.GetContractSwapAggregateInBlockRange)
	router.GET("/contract/swap-aggregate-amount/:codehash/:genesis",
		cache.CacheByRequestURI(store, 60*time.Second), controller.GetContractSwapAggregateAmountInBlockRange)

	router.GET("/ft/info/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllFTInfo)
	router.GET("/ft/codehash/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllFTCodeHash)
	router.GET("/ft/codehash-info/:codehash",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListFTSummary)
	router.GET("/ft/genesis-info/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListFTInfoByGenesis)

	router.GET("/ft/transfer-times/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetFTTransferVolumeInBlockRange)
	router.GET("/ft/owners/:codehash/:genesis",
		cache.CacheByRequestURI(store, 1*time.Second), controller.ListFTOwners)
	router.GET("/ft/summary/:address",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListAllFTSummaryByOwner)

	router.GET("/ft/summary-data/:address",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListAllFTSummaryDataByOwner)

	router.GET("/ft/balance/:codehash/:genesis/:address", controller.GetFTBalanceByOwner) // without cache

	router.GET("/ft/history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetFTHistoryByGenesis)

	router.GET("/ft/income-history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetFTIncomeHistoryByGenesis)

	router.GET("/nft/info/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllNFTInfo)
	router.GET("/nft/codehash/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllNFTCodeHash)
	router.GET("/nft/codehash-info/:codehash",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListNFTSummary)
	router.GET("/nft/genesis-info/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListNFTInfoByGenesis)

	router.GET("/nft/transfer-times/:codehash/:genesis/:tokenid",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetNFTTransferTimesInBlockRange)
	router.GET("/nft/owners/:codehash/:genesis",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListNFTOwners)
	router.GET("/nft/summary/:address",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListAllNFTByOwner)

	// without cache
	router.GET("/nft/detail/:codehash/:genesis/:address", controller.ListNFTCountByOwner)

	router.GET("/nft/history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetNFTHistoryByGenesis)

	router.GET("/address/:address/history",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetHistoryByAddress)

	router.GET("/address/:address/history/tx",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetTxsHistoryByAddress)

	router.GET("/address/:address/contract-history",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetContractHistoryByAddress)

	router.GET("/contract/history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetHistoryByGenesis)

	router.GET("/contract/history/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetAllHistoryByGenesis)

	router.GET("/token/info",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllTokenInfo)

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

	// GC
	go func() {
		for {
			runtime.GC()
			var rtm runtime.MemStats
			runtime.ReadMemStats(&rtm)
			// free memory when large idle
			if rtm.HeapIdle-rtm.HeapReleased > 1*1024*1024*1024 {
				logger.Log.Info("GC",
					zap.String("mAlloc", byteCountBinary(rtm.HeapAlloc)),
					zap.String("mIdle", byteCountBinary(rtm.HeapIdle-rtm.HeapReleased)),
				)
				debug.FreeOSMemory()
			}
			time.Sleep(time.Second * 10)
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

func byteCountBinary(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
