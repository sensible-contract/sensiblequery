package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sensiblequery/logger"
	"sensiblequery/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ybbus/jsonrpc/v2"
	"go.uber.org/zap"
)

var rpcClient jsonrpc.RPCClient

func init() {
	viper.SetConfigFile("conf/chain.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	rpcAddress := viper.GetString("rpc")
	rpcAuth := viper.GetString("rpc_auth")
	rpcClient = jsonrpc.NewClientWithOpts(rpcAddress, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(rpcAuth)),
		},
	})

}

type TxRequest struct {
	TxHex string `json:"txHex"`
}

// LocalPushTx
// @Summary Push Tx to local bitcoind
// @Produce json
// @Param body body TxRequest true "txHex"
// @Success 200 {object} model.Response{data=string} "{"code": 0, "data": "<txid>", "msg": "ok"}"
// @Security BearerAuth
// @Router /local_pushtx [post]
func LocalPushTx(ctx *gin.Context) {
	logger.Log.Info("LocalPushTx enter")

	// check body
	req := TxRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Info("Bind json failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	_, err := hex.DecodeString(req.TxHex)
	if err != nil {
		logger.Log.Info("txRaw invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tx invalid"})
		return
	}

	logger.Log.Info("send", zap.String("rawtx", req.TxHex))
	response, err := rpcClient.Call("sendrawtransaction", []string{req.TxHex})
	if err != nil {
		logger.Log.Info("call failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "rpc failed"})
		return
	}
	logger.Log.Info("Receive remote return", zap.Any("response", response))

	if response.Error != nil {
		ctx.JSON(http.StatusOK, model.Response{
			Code: response.Error.Code,
			Msg:  response.Error.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: response.Result,
	})

}

// WocPushTx
// @Summary Push Tx to woc
// @Produce json
// @Param body body TxRequest true "txHex"
// @Success 200 {object} model.Response{data=string} "{"code": 0, "data": "<txid>", "msg": "ok"}"
// @Security BearerAuth
// @Router /pushtx [post]
func WocPushTx(ctx *gin.Context) {
	logger.Log.Info("WocPushTx enter")

	// check body
	req := TxRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Info("Bind json failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	_, err := hex.DecodeString(req.TxHex)
	if err != nil {
		logger.Log.Info("txRaw invalid", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tx invalid"})
		return
	}

	logger.Log.Info("send", zap.String("rawtx", req.TxHex))

	woc := "https://api.whatsonchain.com/v1/bsv/main/tx/raw"
	if is_testnet != "" {
		woc = "https://api.whatsonchain.com/v1/bsv/test/tx/raw"
	}
	jsonData := fmt.Sprintf(`{"txhex": "%s"}`, req.TxHex)
	resp, err := http.Post(woc, "application/json", bytes.NewBufferString(jsonData))
	if err != nil {
		logger.Log.Info("push tx failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "push tx failed"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	logger.Log.Info("Receive remote return", zap.String("response", string(body)))

	if _, err := hex.DecodeString(string(body)); err != nil {
		ctx.JSON(http.StatusOK, model.Response{
			Code: -1,
			Msg:  string(body),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: string(body),
	})
}

type TxsRequest struct {
	TxsHex []string `json:"txsHex"`
}

// LocalPushTxs
// @Summary Push Tx list to local bitcoind
// @Produce json
// @Param body body TxsRequest true "txsHex"
// @Success 200 {object} model.Response{data=[]string} "{"code": 0, "data": ["<txid>", "<txid>"...], "msg": "ok"}"
// @Security BearerAuth
// @Router /local_pushtxs [post]
func LocalPushTxs(ctx *gin.Context) {
	logger.Log.Info("LocalPushTxs enter")

	// check body
	req := TxsRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Info("Bind json failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	for idx, txHex := range req.TxsHex {
		if len(txHex) == 0 {
			continue
		}
		_, err := hex.DecodeString(txHex)
		if err != nil {
			logger.Log.Info("txRaw invalid", zap.Error(err))
			ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: fmt.Sprintf("tx[%d] invalid", idx)})
			return
		}
	}

	txIdResponse := []interface{}{}
	for _, txHex := range req.TxsHex {
		if len(txHex) == 0 {
			continue
		}

		logger.Log.Info("send", zap.String("rawtx", txHex))
		response, err := rpcClient.Call("sendrawtransaction", []string{txHex})
		if err != nil {
			logger.Log.Info("call failed", zap.Error(err))
			ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "rpc failed", Data: txIdResponse})
			return
		}
		logger.Log.Info("Receive remote return", zap.Any("response", response))

		if response.Error != nil {
			ctx.JSON(http.StatusOK, model.Response{
				Code: response.Error.Code,
				Msg:  response.Error.Message,
				Data: txIdResponse,
			})
			return
		}

		txIdResponse = append(txIdResponse, response.Result)

	}
	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: txIdResponse,
	})

}

// WocPushTxs
// @Summary Push Tx list to woc
// @Produce json
// @Param body body TxsRequest true "txsHex"
// @Success 200 {object} model.Response{data=[]string} "{"code": 0, "data": ["<txid>", "<txid>"...], "msg": "ok"}"
// @Security BearerAuth
// @Router /pushtxs [post]
func WocPushTxs(ctx *gin.Context) {
	logger.Log.Info("WocPushTxs enter")

	// check body
	req := TxsRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Log.Info("Bind json failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	for idx, txHex := range req.TxsHex {
		if len(txHex) == 0 {
			continue
		}
		_, err := hex.DecodeString(txHex)
		if err != nil {
			logger.Log.Info("txRaw invalid", zap.Error(err))
			ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: fmt.Sprintf("tx[%d] invalid", idx)})
			return
		}
	}

	txIdResponse := []interface{}{}
	for _, txHex := range req.TxsHex {
		if len(txHex) == 0 {
			continue
		}

		logger.Log.Info("send", zap.String("rawtx", txHex))

		woc := "https://api.whatsonchain.com/v1/bsv/main/tx/raw"
		if is_testnet != "" {
			woc = "https://api.whatsonchain.com/v1/bsv/test/tx/raw"
		}
		jsonData := fmt.Sprintf(`{"txhex": "%s"}`, txHex)
		resp, err := http.Post(woc, "application/json", bytes.NewBufferString(jsonData))
		if err != nil {
			logger.Log.Info("push tx failed", zap.Error(err))
			ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "push tx failed"})
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		logger.Log.Info("Receive remote return", zap.String("response", string(body)))

		if _, err := hex.DecodeString(string(body)); err != nil {
			ctx.JSON(http.StatusOK, model.Response{
				Code: -1,
				Msg:  string(body),
				Data: txIdResponse,
			})
			return
		}
		txIdResponse = append(txIdResponse, string(body))
	}
	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: txIdResponse,
	})
}

// GetRawMempool
// @Summary GetRawMempool, get txid list in mempool
// @Produce json
// @Success 200 {object} model.Response{data=[]string} "{"code": 0, "data": "[<txid>]", "msg": "ok"}"
// @Security BearerAuth
// @Router /getrawmempool [get]
func GetRawMempool(ctx *gin.Context) {
	logger.Log.Info("GetRawMempool enter")

	response, err := rpcClient.Call("getrawmempool", []string{})
	if err != nil {
		logger.Log.Info("call failed", zap.Error(err))
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "rpc failed"})
		return
	}
	logger.Log.Info("Receive remote return", zap.Any("response", response))

	if response.Error != nil {
		ctx.JSON(http.StatusOK, model.Response{
			Code: response.Error.Code,
			Msg:  response.Error.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: 0,
		Msg:  "ok",
		Data: response.Result,
	})

}
