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
	"sensiblequery/lib/midware"
	"sensiblequery/logger"
	"syscall"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"

	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func SetSwagTitle(title string) func(*ginSwagger.Config) {
	return func(c *ginSwagger.Config) {
		c.Title = title
	}
}

// @title Sensible Query Spec
// @version 2.0
// @description API definition for Sensiblequery  APIs

// @contact.name sensiblequery
// @contact.url https://github.com/sensible-contract/sensiblequery
// @contact.email jiedohh@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger.Log, true))
	router.Use(midware.Metrics())

	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	// go get -u github.com/swaggo/swag/cmd/swag@v1.6.7

	// store := persist.NewMemoryStore(time.Second)
	store := persist.NewRedisStore(rdb.CacheClient)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(basePath+"/swagger/doc.json"),
		SetSwagTitle("Sensible")))

	midware.CreateMetricsEndpoint(router)

	router.GET("/", controller.Satotx)

	router.POST("/pushtx", controller.PushTx)
	router.POST("/pushtxs", controller.PushTxs)

	mainAPI := router.Group("/", midware.VerifyToken())

	mainAPI.GET("/getrawmempool", controller.GetRawMempool)

	mainAPI.GET("/blockchain/info", controller.GetBlockchainInfo)

	mainAPI.GET("/mempool/info", controller.GetMempoolInfo)

	mainAPI.GET("/blocks", controller.GetBlocksByHeightRange)

	mainAPI.GET("/block/id/:blkid", controller.GetBlockById)
	mainAPI.GET("/block/txs/:blkid", controller.GetBlockTxsByBlockId)

	mainAPI.GET("/rawtx/:txid", controller.GetRawTxById)
	mainAPI.GET("/relay/:txid", controller.RelayTxById)

	mainAPI.GET("/tx/:txid", controller.GetTxById)
	mainAPI.GET("/tx/:txid/ins", controller.GetTxInputsByTxId)
	mainAPI.GET("/tx/:txid/outs", controller.GetTxOutputsByTxId)

	mainAPI.GET("/tx/:txid/in/:index", controller.GetTxInputByTxIdAndIdx)
	mainAPI.GET("/tx/:txid/out/:index", controller.GetTxOutputByTxIdAndIdx)

	mainAPI.GET("/tx/:txid/out/:index/spent", controller.GetTxOutputSpentStatusByTxIdAndIdx)

	mainAPI.GET("/address/:address/utxo",
		cache.CacheByRequestURI(store, 1*time.Second), controller.GetUtxoByAddress)
	mainAPI.GET("/address/:address/utxo-data",
		cache.CacheByRequestURI(store, 1*time.Second), controller.GetUtxoDataByAddress)

	mainAPI.GET("/nft/sell/utxo", controller.GetNFTSellUtxo)
	mainAPI.GET("/nft/sell/utxo-by-address/:address", controller.GetNFTSellUtxoByAddress)
	mainAPI.GET("/nft/sell/utxo/:codehash/:genesis", controller.GetNFTSellUtxoByGenesis)
	mainAPI.GET("/nft/sell/utxo-detail/:codehash/:genesis/:token_index", controller.GetNFTSellUtxoDetail)

	mainAPI.GET("/nft/auction/utxo-detail/:codehash/:nftid", controller.GetNFTAuctionUtxoDetail)

	mainAPI.GET("/ft/utxo-data/:codehash/:genesis/:address", controller.GetFTUtxoData)
	mainAPI.GET("/nft/utxo-data/:codehash/:genesis/:address", controller.GetNFTUtxoData)

	mainAPI.GET("/ft/utxo/:codehash/:genesis/:address", controller.GetFTUtxo)
	mainAPI.GET("/nft/utxo/:codehash/:genesis/:address", controller.GetNFTUtxo)
	mainAPI.GET("/nft/utxo-detail/:codehash/:genesis/:token_index", controller.GetNFTUtxoDetailByTokenIndex)
	mainAPI.GET("/nft/utxo-list/:codehash/:genesis",
		cache.CacheByRequestURI(store, 1*time.Second), controller.GetNFTUtxoList)

	mainAPI.GET("/address/:address/balance", controller.GetBalanceByAddress)

	mainAPI.GET("/contract/swap-data/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetContractSwapDataInBlockRange)
	mainAPI.GET("/contract/swap-aggregate/:codehash/:genesis",
		cache.CacheByRequestURI(store, 60*time.Second), controller.GetContractSwapAggregateInBlockRange)
	mainAPI.GET("/contract/swap-aggregate-amount/:codehash/:genesis",
		cache.CacheByRequestURI(store, 60*time.Second), controller.GetContractSwapAggregateAmountInBlockRange)

	mainAPI.GET("/ft/info/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllFTInfo)
	mainAPI.GET("/ft/codehash/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllFTCodeHash)
	mainAPI.GET("/ft/codehash-info/:codehash",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListFTSummary)
	mainAPI.GET("/ft/genesis-info/:codehash/:genesis",
		cache.CacheByRequestURI(store, 60*time.Second), controller.ListFTInfoByGenesis)

	mainAPI.GET("/ft/transfer-times/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetFTTransferVolumeInBlockRange)
	mainAPI.GET("/ft/owners/:codehash/:genesis",
		cache.CacheByRequestURI(store, 1*time.Second), controller.ListFTOwners)
	mainAPI.GET("/ft/summary/:address",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListAllFTSummaryByOwner)

	mainAPI.GET("/ft/summary-data/:address",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListAllFTSummaryDataByOwner)

	mainAPI.GET("/ft/balance/:codehash/:genesis/:address", controller.GetFTBalanceByOwner) // without cache

	mainAPI.GET("/ft/history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetFTHistoryByGenesis)

	mainAPI.GET("/ft/income-history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetFTIncomeHistoryByGenesis)

	mainAPI.GET("/nft/info/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllNFTInfo)
	mainAPI.GET("/nft/codehash/all",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllNFTCodeHash)
	mainAPI.GET("/nft/codehash-info/:codehash",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListNFTSummary)
	mainAPI.GET("/nft/genesis-info/:codehash/:genesis",
		cache.CacheByRequestURI(store, 60*time.Second), controller.ListNFTInfoByGenesis)

	mainAPI.GET("/nft/transfer-times/:codehash/:genesis/:tokenid",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetNFTTransferTimesInBlockRange)
	mainAPI.GET("/nft/owners/:codehash/:genesis",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListNFTOwners)
	mainAPI.GET("/nft/summary/:address",
		cache.CacheByRequestURI(store, 2*time.Second), controller.ListAllNFTByOwner)

	// without cache
	mainAPI.GET("/nft/detail/:codehash/:genesis/:address", controller.ListNFTCountByOwner)

	mainAPI.GET("/nft/history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetNFTHistoryByGenesis)

	mainAPI.GET("/address/:address/history",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetHistoryByAddress)

	mainAPI.GET("/address/:address/history/tx",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetTxsHistoryByAddress)

	mainAPI.GET("/address/:address/contract-history",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetContractHistoryByAddress)

	mainAPI.GET("/contract/history/:codehash/:genesis/:address",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetHistoryByGenesis)

	mainAPI.GET("/contract/history/:codehash/:genesis",
		cache.CacheByRequestURI(store, 10*time.Second), controller.GetAllHistoryByGenesis)

	mainAPI.GET("/token/info",
		cache.CacheByRequestURI(store, 10*time.Second), controller.ListAllTokenInfo)

	heightAPI := router.Group("/height/:height", midware.VerifyToken())
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
