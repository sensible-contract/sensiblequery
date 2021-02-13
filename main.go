package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"satoblock/controller"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	// 0.0.0.0:8000
	listen_address = os.Getenv("LISTEN")
)

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

	router.GET("/", controller.Satotx)
	router.GET("/blockchain/info", controller.GetBlockchainInfo)

	router.GET("/blocks/:start/:end", controller.GetBlocksByHeightRange)

	router.GET("/block/id/:blkid", controller.GetBlockById)
	router.GET("/block/txs/:blkid", controller.GetBlockTxsByBlockId)

	router.GET("/tx/:txid", controller.GetTxById)
	router.GET("/tx/:txid/ins", controller.GetTxInputsByTxId)
	router.GET("/tx/:txid/outs", controller.GetTxOutputsByTxId)

	router.GET("/tx/:txid/in/:index", controller.GetTxInputByTxIdAndIdx)
	router.GET("/tx/:txid/out/:index", controller.GetTxOutputByTxIdAndIdx)

	router.GET("/tx/:txid/out/:index/spent", controller.GetTxOutputSpentStatusByTxIdAndIdx)

	router.GET("/address/:address/utxo", controller.GetUtxoByAddress)
	router.GET("/genesis/:genesis/utxo", controller.GetUtxoByGenesis)

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

	beforeHeightAPI := router.Group("/before/:height")
	{
		beforeHeightAPI.GET("/tx/:txid", controller.GetTxByIdBeforeHeight)
		beforeHeightAPI.GET("/tx/:txid/ins", controller.GetTxInputsByTxIdBeforeHeight)
		beforeHeightAPI.GET("/tx/:txid/outs", controller.GetTxOutputsByTxIdBeforeHeight)
	}

	afterHeightAPI := router.Group("/after/:height")
	{
		afterHeightAPI.GET("/tx/:txid", controller.GetTxByIdAfterHeight)
		// afterHeightAPI.GET("/tx/:txid/ins", controller.GetTxInputsByTxIdAfterHeight)
		// afterHeightAPI.GET("/tx/:txid/outs", controller.GetTxOutputsByTxIdAfterHeight)
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

	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
