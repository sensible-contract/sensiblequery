package service

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"satosensible/lib/blkparser"
	"satosensible/lib/utils"
	"satosensible/logger"
	"satosensible/model"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

func mergeUtxoByKeys(addressUtxoConfirmed, addressUtxoSpentUnconfirmed, oldUtxoKey, newUtxoKey, finalKey string) (err error) {
	// 注意这里查询需要原子化，可使用pipeline
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

	return nil
}

////////////////
func getNFTSellUtxoFromRedis(utxoOutpoints []string) (nftSellsRsp []*model.NFTSellResp, err error) {
	logger.Log.Info("getNFTSellUtxoFromRedis redis", zap.Int("nUTXO", len(utxoOutpoints)))
	nftSellsRsp = make([]*model.NFTSellResp, 0)
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
		txout := &model.TxoData{}
		txout.Unmarshal([]byte(res))

		// 补充数据
		txout.UTxid = []byte(outpoint[:32])                            // 32
		txout.Vout = binary.LittleEndian.Uint32([]byte(outpoint[32:])) // 4
		txout.ScriptType = scriptDecoder.GetLockingScriptType(txout.Script)

		txo := scriptDecoder.ExtractPkScriptForTxo(txout.Script, txout.ScriptType)

		nftSellsRsp = append(nftSellsRsp, &model.NFTSellResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txo.AddressPkh, utils.PubKeyHashAddrID),
			Satoshi: int(txout.Satoshi),

			TokenIndex:      strconv.FormatUint(txo.TokenIndex, 10),
			MetaTxIdHex:     hex.EncodeToString(txo.MetaTxId),
			MetaOutputIndex: int(txo.MetaOutputIndex),
			CodeHashHex:     hex.EncodeToString(txo.CodeHash),
			GenesisHex:      hex.EncodeToString(txo.GenesisId),
			SensibleIdHex:   hex.EncodeToString(txo.SensibleId),
			Height:          int(txout.BlockHeight),
			Idx:             int(txout.TxIdx),
		})
	}

	return nftSellsRsp, nil
}

////////////////
func GetNFTSellUtxo(cursor, size int) (nftSellsRsp []*model.NFTSellResp, err error) {
	addressUtxoConfirmed := "{sut}"
	addressUtxoSpentUnconfirmed := "mp:s:{sut}"
	oldUtxoKey := "mp:t:{sut}"
	newUtxoKey := "mp:{sut}"
	finalKey := "mp:z:{sut}"

	if err := mergeUtxoByKeys(addressUtxoConfirmed, addressUtxoSpentUnconfirmed, oldUtxoKey, newUtxoKey, finalKey); err != nil {
		return nil, err
	}

	//
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetNFTSellUtxo redis failed", zap.Error(err))
		return
	}
	return getNFTSellUtxoFromRedis(utxoOutpoints)
}

//////////////// address
func GetNFTSellUtxoByAddress(cursor, size int, addressPkh []byte) (nftSellsRsp []*model.NFTSellResp, err error) {
	addressKey := string(addressPkh) + "}"

	addressUtxoConfirmed := "{suta" + addressKey
	addressUtxoSpentUnconfirmed := "mp:s:{suta" + addressKey
	oldUtxoKey := "mp:t:{suta" + addressKey
	newUtxoKey := "mp:{suta" + addressKey
	finalKey := "mp:z:{suta" + addressKey

	if err := mergeUtxoByKeys(addressUtxoConfirmed, addressUtxoSpentUnconfirmed, oldUtxoKey, newUtxoKey, finalKey); err != nil {
		return nil, err
	}

	//
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetNFTSellUtxoByAddress redis failed", zap.Error(err))
		return
	}
	return getNFTSellUtxoFromRedis(utxoOutpoints)
}

//////////////// genesisId
func GetNFTSellUtxoByGenesis(cursor, size int, codeHash, genesisId []byte) (nftSellsRsp []*model.NFTSellResp, err error) {
	genesisKey := string(genesisId) + string(codeHash) + "}"

	addressUtxoConfirmed := "{sutc" + genesisKey
	addressUtxoSpentUnconfirmed := "mp:s:{sutc" + genesisKey
	oldUtxoKey := "mp:t:{sutc" + genesisKey
	newUtxoKey := "mp:{sutc" + genesisKey
	finalKey := "mp:z:{sutc" + genesisKey

	if err := mergeUtxoByKeys(addressUtxoConfirmed, addressUtxoSpentUnconfirmed, oldUtxoKey, newUtxoKey, finalKey); err != nil {
		return nil, err
	}

	//
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetNFTSellUtxoByGenesis redis failed", zap.Error(err))
		return
	}
	return getNFTSellUtxoFromRedis(utxoOutpoints)
}

//////////////// address utxo
func GetNFTSellUtxoByTokenIndexMerge(codeHash, genesisId []byte, tokenIndex string) (nftSellsRsp *model.NFTSellResp, err error) {
	key := "mp:{suic" + string(genesisId) + string(codeHash) + "}"
	resp, err := GetNFTSellUtxoByTokenIndex(key, tokenIndex)
	if err == nil {
		return resp, nil
	}

	key = "{suic" + string(genesisId) + string(codeHash) + "}"
	return GetNFTSellUtxoByTokenIndex(key, tokenIndex)
}

func GetNFTSellUtxoByTokenIndex(key string, tokenIndex string) (nftSellsRsp *model.NFTSellResp, err error) {
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
	result, err := getNFTSellUtxoFromRedis(utxoOutpoints)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("not exist")
	}
	return result[0], nil
}
