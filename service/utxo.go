package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"satosensible/lib/blkparser"
	"satosensible/lib/utils"
	"satosensible/logger"
	"satosensible/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rdb      *redis.Client
	rdbBlock *redis.Client
	ctx      = context.Background()
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

	address := viper.GetString("address")
	password := viper.GetString("password")
	databaseBlock := viper.GetInt("database_block")
	database := viper.GetInt("database")
	dialTimeout := viper.GetDuration("dialTimeout")
	readTimeout := viper.GetDuration("readTimeout")
	writeTimeout := viper.GetDuration("writeTimeout")
	poolSize := viper.GetInt("poolSize")
	rdb = redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           database,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     poolSize,
	})

	rdbBlock = redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           databaseBlock,
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

	balance, err := rdb.ZScore(ctx, "balance", string(addressPkh)).Result()
	if err == redis.Nil {
		balance = 0
	} else if err != nil {
		logger.Log.Info("GetBalanceByAddress redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("GetBalanceByAddress", zap.Float64("balance", balance))
	balanceRsp.Satoshi = int(balance)

	// 待确认余额
	mpBalance, err := rdb.ZScore(ctx, "mp:balance", string(addressPkh)).Result()
	if err == redis.Nil {
		mpBalance = 0
	} else if err != nil {
		logger.Log.Info("GetBalanceByAddress redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("GetBalanceByAddress", zap.Float64("pending", mpBalance))
	balanceRsp.PendingSatoshi = int(mpBalance)

	return balanceRsp, nil
}

//////////////// address utxo
func GetUtxoByTokenId(codeHash, genesisId []byte, tokenId string) (txOutsRsp *model.TxOutResp, err error) {
	key := "mp:nd" + string(codeHash) + string(genesisId)
	resp, err := GetNFTUtxoByTokenId(key, tokenId)
	if err == nil {
		return resp, nil
	}

	key = "nd" + string(codeHash) + string(genesisId)
	return GetNFTUtxoByTokenId(key, tokenId)
}

func GetNFTUtxoByTokenId(key string, tokenId string) (txOutsRsp *model.TxOutResp, err error) {
	op := &redis.ZRangeBy{
		Min:    tokenId, // 最小分数
		Max:    tokenId, // 最大分数
		Offset: 0,       // 类似sql的limit, 表示开始偏移量
		Count:  1,       // 一次返回多少数据
	}
	utxoOutpoints, err := rdb.ZRangeByScore(ctx, key, op).Result()
	if err != nil {
		logger.Log.Info("GetUtxoByTokenId redis failed", zap.Error(err))
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
	addressKey := string(codeHash) + string(genesisId) + string(addressPkh)
	addressUtxoConfirmed := ""
	addressUtxoSpentUnconfirmed := ""
	oldUtxoKey := ""
	newUtxoKey := ""

	if isNFT {
		addressUtxoConfirmed = "nu" + addressKey
		addressUtxoSpentUnconfirmed = "mp:s:nu" + addressKey
		oldUtxoKey = "mp:t:nu" + addressKey
		newUtxoKey = "mp:nu" + addressKey
		finalKey = "mp:z:nu" + addressKey
	} else {
		addressUtxoConfirmed = "fu" + addressKey
		addressUtxoSpentUnconfirmed = "mp:s:fu" + addressKey
		oldUtxoKey = "mp:t:fu" + addressKey
		newUtxoKey = "mp:fu" + addressKey
		finalKey = "mp:z:fu" + addressKey
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
	pipe := rdbBlock.Pipeline()

	outpointsCmd := make([]*redis.StringCmd, 0)
	for _, outpoint := range utxoOutpoints {
		outpointsCmd = append(outpointsCmd, pipe.Get(ctx, outpoint))
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
		txout := &model.TxoData{}
		txout.Unmarshal([]byte(res))

		// 补充数据
		txout.UTxid = []byte(outpoint[:32])                            // 32
		txout.Vout = binary.LittleEndian.Uint32([]byte(outpoint[32:])) // 4
		txout.ScriptType = scriptDecoder.GetLockingScriptType(txout.Script)

		txo := scriptDecoder.ExtractPkScriptForTxo(txout.Script, txout.ScriptType)
		tokenId := ""
		if len(txo.GenesisId) >= 20 {
			if txo.CodeType == scriptDecoder.CodeType_NFT {
				tokenId = strconv.FormatUint(txo.TokenIdx, 10)
			} else if txo.CodeType == scriptDecoder.CodeType_FT || txo.CodeType == scriptDecoder.CodeType_UNIQUE {
				tokenId = hex.EncodeToString(txo.GenesisId)
			}
		}

		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txo.AddressPkh, utils.PubKeyHashAddrID),
			Satoshi: int(txout.Satoshi),

			IsNFT:         (txo.CodeType == scriptDecoder.CodeType_NFT),
			CodeType:      int(txo.CodeType),
			TokenId:       tokenId,
			MetaTxId:      hex.EncodeToString(txo.MetaTxId),
			TokenName:     txo.Name,
			TokenSymbol:   txo.Symbol,
			TokenAmount:   strconv.FormatUint(txo.Amount, 10),
			TokenDecimal:  int(txo.Decimal),
			CodeHashHex:   hex.EncodeToString(txo.CodeHash),
			GenesisHex:    hex.EncodeToString(txo.GenesisId),
			SensibleIdHex: hex.EncodeToString(txo.SensibleId),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			// ScriptPkHex:   hex.EncodeToString(txout.Script),
			Height: int(txout.BlockHeight),
			Idx:    int(txout.TxIdx),
		})
	}

	return txOutsRsp, nil
}

//////////////// address utxo
func GetUtxoByAddress(cursor, size int, addressPkh []byte) (txOutsRsp []*model.TxStandardOutResp, err error) {
	logger.Log.Info("GetUtxoByAddress", zap.String("addressHex", hex.EncodeToString(addressPkh)))

	addressUtxoConfirmed := "au" + string(addressPkh)
	addressUtxoSpentUnconfirmed := "mp:s:au" + string(addressPkh)
	oldUtxoKey := "mp:t:au" + string(addressPkh)
	newUtxoKey := "mp:au" + string(addressPkh)
	finalKey := "mp:z:au" + string(addressPkh)

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
	logger.Log.Info("getUtxoFromRedis", zap.Int("nOutpoints", len(utxoOutpoints)))
	txOutsRsp = make([]*model.TxStandardOutResp, 0)
	pipe := rdbBlock.Pipeline()

	outpointsCmd := make([]*redis.StringCmd, 0)
	for _, outpoint := range utxoOutpoints {
		outpointsCmd = append(outpointsCmd, pipe.Get(ctx, outpoint))
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
		txout := &model.TxoData{}
		txout.Unmarshal([]byte(res))

		// 补充数据
		txout.UTxid = []byte(outpoint[:32])                            // 32
		txout.Vout = binary.LittleEndian.Uint32([]byte(outpoint[32:])) // 4
		txout.ScriptType = scriptDecoder.GetLockingScriptType(txout.Script)

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
