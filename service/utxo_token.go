package service

import (
	"encoding/hex"
	"errors"
	"sensiblequery/logger"
	"sensiblequery/model"
	"sort"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

//////////////// address NFT utxo
func GetUtxoByTokenIndex(codeHash, genesisId []byte, tokenIndex string) (txOutsRsp *model.TxOutResp, err error) {
	key := "mp:nd" + string(codeHash) + string(genesisId)
	resp, err := getNFTUtxoByTokenIndex(key, tokenIndex)
	if err == nil {
		return resp, nil
	}

	key = "nd" + string(codeHash) + string(genesisId)
	return getNFTUtxoByTokenIndex(key, tokenIndex)
}

func getNFTUtxoByTokenIndex(key string, tokenIndex string) (txOutsRsp *model.TxOutResp, err error) {
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
		txOutRsp.ScriptPkHex = ""
		txOutsRsp = append(txOutsRsp, txOutRsp)
	}

	return txOutsRsp, nil
}

//////////////// address utxo
func GetUtxoByCodeHashGenesisAddress(cursor, size int, codeHash, genesisId, addressPkh []byte, key string) (
	txOutsRsp []*model.TxOutResp, total, totalConf, totalUnconf, totalUnconfSpend int, err error) {
	logger.Log.Info("GetUtxoByCodeHashGenesisAddress",
		zap.String("codehash", hex.EncodeToString(codeHash)),
		zap.String("genesis", hex.EncodeToString(genesisId)),
		zap.String("addressHex", hex.EncodeToString(addressPkh)),
	)

	utxoOutpoints, total, totalConf, totalUnconf, totalUnconfSpend, err := GetUtxoOutpointsByAddress(cursor, size, codeHash, genesisId, addressPkh, key)
	if err != nil {
		return
	}

	txOutsRsp, err = getUtxoFromRedis(utxoOutpoints)
	return txOutsRsp, total, totalConf, totalUnconf, totalUnconfSpend, err
}

//////////////// list NFT utxo
func GetNFTUtxoByTokenIndexRange(cursor, size int, codeHash, genesisId []byte) (
	txOutsRsp []*model.TxOutResp, total, totalConf, totalUnconf int, err error) {
	newUtxoKey := "mp:nd" + string(codeHash) + string(genesisId)

	// unconfirmed count
	newUtxoNum, err := rdb.ZCard(ctx, newUtxoKey).Result()
	if err != nil {
		logger.Log.Info("get newUtxoNum from redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("newUtxoNum", zap.Int64("n", newUtxoNum))

	utxoKeyConfirmed := "nd" + string(codeHash) + string(genesisId)
	// confirmed count
	utxoConfirmedNum, err := rdb.ZCard(ctx, utxoKeyConfirmed).Result()
	if err != nil {
		logger.Log.Info("get utxoConfirmedNum from redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("utxoConfirmedNum", zap.Int64("n", utxoConfirmedNum))

	totalConf = int(utxoConfirmedNum)
	totalUnconf = int(newUtxoNum)
	total = totalConf + totalUnconf

	// mempool
	op := &redis.ZRangeBy{
		Min: strconv.Itoa(cursor),            // 最小分数
		Max: strconv.Itoa(cursor + size - 1), // 最大分数
		// Offset: 0,                               // 类似sql的limit, 表示开始偏移量
		// Count:  1,                               // 一次返回多少数据
	}
	utxoOutpointsUnconfirmed, err := rdb.ZRangeByScoreWithScores(ctx, newUtxoKey, op).Result()
	if err != nil {
		logger.Log.Info("GetNFTUtxoByTokenIndexRange redis failed", zap.Error(err))
		return
	}

	utxoOutpointsWithScore, err := rdb.ZRangeByScoreWithScores(ctx, utxoKeyConfirmed, op).Result()
	if err != nil {
		logger.Log.Info("GetNFTUtxoByTokenIndexRange redis failed", zap.Error(err))
		return
	}

	utxoOutpointsMap := make(map[int]string, size)
	for _, out := range utxoOutpointsWithScore {
		utxoOutpointsMap[int(out.Score)] = out.Member.(string)
	}
	for _, out := range utxoOutpointsUnconfirmed {
		utxoOutpointsMap[int(out.Score)] = out.Member.(string)
	}

	var utxoOutpoints []string
	// sort
	var keys []int
	for k := range utxoOutpointsMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		utxoOutpoints = append(utxoOutpoints, utxoOutpointsMap[k])
	}

	// get utxo data
	result, err := getUtxoFromRedis(utxoOutpoints)
	if err != nil {
		return
	}
	if len(result) == 0 {
		err = errors.New("not exist")
		return
	}
	return result, total, totalConf, totalUnconf, nil
}
