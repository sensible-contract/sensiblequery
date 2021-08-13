package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"satosensible/lib/blkparser"
	"satosensible/lib/utils"
	"satosensible/logger"
	"satosensible/model"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rdb redis.UniversalClient
	ctx = context.Background()
)

func init() {
	viper.SetConfigFile("conf/redis.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	addrs := viper.GetStringSlice("addrs")
	password := viper.GetString("password")
	database := viper.GetInt("database")
	dialTimeout := viper.GetDuration("dialTimeout")
	readTimeout := viper.GetDuration("readTimeout")
	writeTimeout := viper.GetDuration("writeTimeout")
	poolSize := viper.GetInt("poolSize")
	rdb = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        addrs,
		Password:     password,
		DB:           database,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     poolSize,
	})
}

func GetBalanceByAddress(addressPkh []byte) (balanceRsp *model.BalanceResp, err error) {
	balanceRsp = &model.BalanceResp{
		Address: utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrID),
	}

	balance, err := rdb.Get(ctx, "bl"+string(addressPkh)).Int()
	if err == redis.Nil {
		balance = 0
	} else if err != nil {
		logger.Log.Info("GetBalanceByAddress redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("GetBalanceByAddress", zap.Int("balance", balance))
	balanceRsp.Satoshi = balance

	// 待确认余额
	mpBalance, err := rdb.Get(ctx, "mp:bl"+string(addressPkh)).Int()
	if err == redis.Nil {
		mpBalance = 0
	} else if err != nil {
		logger.Log.Info("GetBalanceByAddress redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("GetBalanceByAddress", zap.Int("pending", mpBalance))
	balanceRsp.PendingSatoshi = mpBalance

	return balanceRsp, nil
}

//////////////// address utxo
func GetUtxoByTokenIndex(codeHash, genesisId []byte, tokenIndex string) (txOutsRsp *model.TxOutResp, err error) {
	key := "mp:nd" + string(codeHash) + string(genesisId)
	resp, err := GetNFTUtxoByTokenIndex(key, tokenIndex)
	if err == nil {
		return resp, nil
	}

	key = "nd" + string(codeHash) + string(genesisId)
	return GetNFTUtxoByTokenIndex(key, tokenIndex)
}

func GetNFTUtxoByTokenIndex(key string, tokenIndex string) (txOutsRsp *model.TxOutResp, err error) {
	op := &redis.ZRangeBy{
		Min:    tokenIndex, // 最小分数
		Max:    tokenIndex, // 最大分数
		Offset: 0,          // 类似sql的limit, 表示开始偏移量
		Count:  1,          // 一次返回多少数据
	}
	utxoOutpoints, err := rdb.ZRangeByScore(ctx, key, op).Result()
	if err != nil {
		logger.Log.Info("GetUtxoByTokenIndex redis failed", zap.Error(err))
		return
	}
	result, err := getUtxoFromRedis(utxoOutpoints)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("not exist")
	}
	return result[0], nil
}

//////////////// merge ft utxo
func mergeUtxoByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte, isNFT bool) (finalKey string, err error) {
	// 注意这里查询需要原子化，可使用pipeline
	addressKey := string(addressPkh) + "}" + string(codeHash) + string(genesisId)
	addressUtxoConfirmed := ""
	addressUtxoSpentUnconfirmed := ""
	oldUtxoKey := ""
	newUtxoKey := ""

	if isNFT {
		addressUtxoConfirmed = "{nu" + addressKey
		addressUtxoSpentUnconfirmed = "mp:s:{nu" + addressKey
		oldUtxoKey = "mp:t:{nu" + addressKey
		newUtxoKey = "mp:{nu" + addressKey
		finalKey = "mp:z:{nu" + addressKey
	} else {
		addressUtxoConfirmed = "{fu" + addressKey
		addressUtxoSpentUnconfirmed = "mp:s:{fu" + addressKey
		oldUtxoKey = "mp:t:{fu" + addressKey
		newUtxoKey = "mp:{fu" + addressKey
		finalKey = "mp:z:{fu" + addressKey
	}

	tmpZs := &redis.ZStore{
		Keys: []string{
			addressUtxoConfirmed, addressUtxoSpentUnconfirmed,
		},
	}
	nDiff, err := rdb.ZDiffStore(ctx, oldUtxoKey, tmpZs).Result()
	if err != nil {
		logger.Log.Info("ZDiffStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZDiffStore", zap.Int64("n", nDiff))

	finalZs := &redis.ZStore{
		Keys: []string{
			oldUtxoKey, newUtxoKey,
		},
	}
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZDiffStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))

	return finalKey, nil
}

//////////////// genesisId
func GetUtxoByCodeHashGenesisAddress(cursor, size int, codeHash, genesisId, addressPkh []byte, isNFT bool) (txOutsRsp []*model.TxOutResp, err error) {
	finalKey, err := mergeUtxoByCodeHashGenesisAddress(codeHash, genesisId, addressPkh, isNFT)
	if err != nil {
		return nil, err
	}

	//
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetUtxoByCodeHashGenesisAddress redis failed", zap.Error(err))
		return
	}
	return getUtxoFromRedis(utxoOutpoints)
}

////////////////
func getUtxoFromRedis(utxoOutpoints []string) (txOutsRsp []*model.TxOutResp, err error) {
	logger.Log.Info("getUtxoFromRedis redis", zap.Int("nUTXO", len(utxoOutpoints)))
	txOutsRsp = make([]*model.TxOutResp, 0)
	pipe := rdb.Pipeline()

	outpointsCmd := make([]*redis.StringCmd, 0)
	for _, outpoint := range utxoOutpoints {
		outpointsCmd = append(outpointsCmd, pipe.Get(ctx, "u"+outpoint))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for outpointIdx, data := range outpointsCmd {
		outpoint := utxoOutpoints[outpointIdx]
		res, err := data.Result()
		if err == redis.Nil {
			logger.Log.Info("redis not found", zap.String("outpoint", hex.EncodeToString([]byte(outpoint))))
			continue
		} else if err != nil {
			panic(err)
		}

		txout := model.NewTxoData([]byte(outpoint), []byte(res))
		txOutDO := model.TxOutDO{
			Height:     txout.BlockHeight,
			Idx:        uint32(txout.TxIdx),
			ScriptPk:   txout.Script,
			Satoshi:    txout.Satoshi,
			TxId:       txout.UTxid, // 32
			Vout:       txout.Vout,  // 4
			ScriptType: txout.ScriptType,
		}
		txOutRsp := getTxOutputRespFromDo(&txOutDO)
		txOutsRsp = append(txOutsRsp, txOutRsp)
	}

	return txOutsRsp, nil
}

//////////////// address utxo
func GetUtxoByAddress(cursor, size int, addressPkh []byte) (txOutsRsp []*model.TxStandardOutResp, err error) {
	logger.Log.Info("GetUtxoByAddress", zap.String("addressHex", hex.EncodeToString(addressPkh)))

	addressUtxoConfirmed := "{au" + string(addressPkh) + "}"
	addressUtxoSpentUnconfirmed := "mp:s:{au" + string(addressPkh) + "}"
	oldUtxoKey := "mp:t:{au" + string(addressPkh) + "}"
	newUtxoKey := "mp:{au" + string(addressPkh) + "}"
	finalKey := "mp:z:{au" + string(addressPkh) + "}"

	tmpZs := &redis.ZStore{
		Keys: []string{
			addressUtxoConfirmed, addressUtxoSpentUnconfirmed,
		},
	}
	nDiff, err := rdb.ZDiffStore(ctx, oldUtxoKey, tmpZs).Result()
	if err != nil {
		logger.Log.Info("ZDiffStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZDiffStore", zap.Int64("n", nDiff))

	finalZs := &redis.ZStore{
		Keys: []string{
			oldUtxoKey, newUtxoKey,
		},
	}
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZUnionStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetUtxoByAddress redis failed", zap.Error(err))
		return
	}
	return getNonTokenUtxoFromRedis(utxoOutpoints)
}

////////////////
func getNonTokenUtxoFromRedis(utxoOutpoints []string) (txOutsRsp []*model.TxStandardOutResp, err error) {
	logger.Log.Info("getNonTokenUtxoFromRedis", zap.Int("nOutpoints", len(utxoOutpoints)))
	txOutsRsp = make([]*model.TxStandardOutResp, 0)
	pipe := rdb.Pipeline()

	outpointsCmd := make([]*redis.StringCmd, 0)
	for _, outpoint := range utxoOutpoints {
		outpointsCmd = append(outpointsCmd, pipe.Get(ctx, "u"+outpoint))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for outpointIdx, data := range outpointsCmd {
		outpoint := utxoOutpoints[outpointIdx]
		res, err := data.Result()
		if err == redis.Nil {
			logger.Log.Info("redis not found", zap.String("outpoint", hex.EncodeToString([]byte(outpoint))))
			continue
		} else if err != nil {
			panic(err)
		}

		txout := model.NewTxoData([]byte(outpoint), []byte(res))
		txOutsRsp = append(txOutsRsp, &model.TxStandardOutResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Satoshi: int(txout.Satoshi),

			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			// ScriptPkHex:   hex.EncodeToString(txout.Script),
			Height: int(txout.BlockHeight),
			Idx:    int(txout.TxIdx),
		})
	}

	return txOutsRsp, nil
}
