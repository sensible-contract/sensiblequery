package controller

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"satosensible/model"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ybbus/jsonrpc/v2"
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

// Pushtx
// @Summary Push Tx
// @Produce json
// @Param body body TxRequest true "txHex"
// @Success 200 {object} model.Response{data=string} "{"code": 0, "data": "<txid>", "msg": "ok"}"
// @Router /pushtx [post]
func PushTx(ctx *gin.Context) {
	log.Printf("PushTx enter")

	// check body
	req := TxRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("Bind json failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	_, err := hex.DecodeString(req.TxHex)
	if err != nil {
		log.Printf("txRaw invalid: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "tx invalid"})
		return
	}

	log.Printf("rawtx: %s:", req.TxHex)
	response, err := rpcClient.Call("sendrawtransaction", []string{req.TxHex})
	if err != nil {
		log.Println("call failed:", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "rpc failed"})
		return
	}
	log.Println("Receive remote return:", response)

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

type TxsRequest struct {
	TxsHex []string `json:"txsHex"`
}

// Pushtxs
// @Summary Push Tx list
// @Produce json
// @Param body body TxsRequest true "txsHex"
// @Success 200 {object} model.Response{data=[]string} "{"code": 0, "data": ["<txid>", "<txid>"...], "msg": "ok"}"
// @Router /pushtxs [post]
func PushTxs(ctx *gin.Context) {
	log.Printf("PushTxs enter")

	// check body
	req := TxsRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("Bind json failed: %v", err)
		ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "json error"})
		return
	}

	for idx, txHex := range req.TxsHex {
		if len(txHex) == 0 {
			continue
		}
		_, err := hex.DecodeString(txHex)
		if err != nil {
			log.Printf("txRaw invalid: %v", err)
			ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: fmt.Sprintf("tx[%d] invalid", idx)})
			return
		}
	}

	txIdResponse := []interface{}{}
	for _, txHex := range req.TxsHex {
		if len(txHex) == 0 {
			continue
		}

		log.Printf("rawtx: %s:", txHex)
		response, err := rpcClient.Call("sendrawtransaction", []string{txHex})
		if err != nil {
			log.Println("call failed:", err)
			ctx.JSON(http.StatusOK, model.Response{Code: -1, Msg: "rpc failed", Data: txIdResponse})
			return
		}
		log.Println("Receive remote return:", response)

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
