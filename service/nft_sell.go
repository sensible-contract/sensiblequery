package service

import (
	"encoding/hex"
	"sensiblequery/dao/rdb"
	"sensiblequery/lib/blkparser"
	"sensiblequery/lib/utils"
	"sensiblequery/logger"
	"sensiblequery/model"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

func mergeUtxoByKeys(addressUtxoConfirmed, addressUtxoSpentUnconfirmed, oldUtxoKey, newUtxoKey, finalKey string) (err error) {
	// 注意这里查询需要原子化，可使用pipeline
	nDiff, err := rdb.BizClient.ZDiffStore(ctx, oldUtxoKey, addressUtxoConfirmed, addressUtxoSpentUnconfirmed).Result()
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
	nUnion, err := rdb.BizClient.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZUnionStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))

	return nil
}

////////////////
func getNFTSellUtxoFromRedis(utxoOutpoints []string) (nftSellsRsp []*model.NFTSellResp, err error) {
	logger.Log.Info("getNFTSellUtxoFromRedis redis", zap.Int("nUTXO", len(utxoOutpoints)))
	nftSellsRsp = make([]*model.NFTSellResp, 0)
	pipe := rdb.RdbUtxoClient.Pipeline()

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
		nftSellRsp := &model.NFTSellResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Satoshi: int(txout.Satoshi),
			Height:  int(txout.BlockHeight),
			Idx:     int(txout.TxIdx),
		}
		txo := scriptDecoder.ExtractPkScriptForTxo(txout.PkScript, txout.ScriptType)
		if txo.CodeType != scriptDecoder.CodeType_NONE && txo.CodeType != scriptDecoder.CodeType_SENSIBLE {
			nftSellRsp.CodeHashHex = hex.EncodeToString(txo.CodeHash[:])
			nftSellRsp.GenesisHex = hex.EncodeToString(txo.GenesisId[:txo.GenesisIdLen])
		}
		if txo.CodeType == scriptDecoder.CodeType_NFT_SELL {
			nftSellRsp.Address = utils.EncodeAddress(txo.AddressPkh[:], utils.PubKeyHashAddrID)
			nftSellRsp.TokenIndex = strconv.FormatUint(txo.NFTSell.TokenIndex, 10)
			nftSellRsp.Price = int(txo.NFTSell.Price)

			// 设置准备状态
			contractHashAsAddressPkh := blkparser.GetHash160(txout.PkScript)
			countRsp, err := GetNFTCountByCodeHashGenesisAddress(txo.CodeHash[:], txo.GenesisId[:txo.GenesisIdLen], contractHashAsAddressPkh)
			if err == nil && countRsp.Count+countRsp.PendingCount > 0 {
				nftSellRsp.IsReady = true
			}
		}

		nftSellsRsp = append(nftSellsRsp, nftSellRsp)
	}

	getNFTMetaInfoForSell(nftSellsRsp)

	return nftSellsRsp, nil
}

func getNFTMetaInfoForSell(nftSellsRsp []*model.NFTSellResp) {
	pipe := rdb.BizClient.Pipeline()
	nftInfoCmds := make([]*redis.StringStringMapCmd, 0)
	for _, nft := range nftSellsRsp {
		// nftinfo of each token
		key, _ := hex.DecodeString(nft.CodeHashHex + nft.GenesisHex)
		nftInfoCmds = append(nftInfoCmds, pipe.HGetAll(ctx, "nI"+string(key)+nft.TokenIndex))
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}
	for idx, nft := range nftSellsRsp {
		// sensible/supply
		if nftinfo, err := nftInfoCmds[idx].Result(); err == nil {
			supply, _ := strconv.Atoi(nftinfo["supply"])
			nft.Supply = supply
			metavout, _ := strconv.Atoi(nftinfo["metavout"])
			nft.MetaOutputIndex = metavout
			nft.MetaTxIdHex = hex.EncodeToString([]byte(nftinfo["metatxid"]))
			nft.SensibleIdHex = hex.EncodeToString([]byte(nftinfo["sensibleid"]))
		} else if err == redis.Nil {
			logger.Log.Info("getNFTMetaInfoForSell nftinfo not found")
		} else if err != nil {
			logger.Log.Info("getNFTMetaInfoForSell redis failed", zap.Error(err))
		}
	}
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

	utxoOutpoints, err := rdb.BizClient.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
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

	utxoOutpoints, err := rdb.BizClient.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
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

	utxoOutpoints, err := rdb.BizClient.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetNFTSellUtxoByGenesis redis failed", zap.Error(err))
		return
	}
	return getNFTSellUtxoFromRedis(utxoOutpoints)
}

//////////////// address utxo
func GetNFTSellUtxoByTokenIndexMerge(codeHash, genesisId []byte, tokenIndex string, isReadyOnly bool) (nftSellsRsp []*model.NFTSellResp, err error) {
	// fixme: 可能被恶意创建sell utxo
	key := "mp:{suic" + string(genesisId) + string(codeHash) + "}"
	respMp, err := GetNFTSellUtxoByTokenIndex(key, tokenIndex)
	if err != nil {
		return nil, err
	}

	key = "{suic" + string(genesisId) + string(codeHash) + "}"
	resp, err := GetNFTSellUtxoByTokenIndex(key, tokenIndex)
	if err != nil {
		return nil, err
	}

	nftSellsRsp = make([]*model.NFTSellResp, 0)
	for _, data := range respMp {
		if isReadyOnly && !data.IsReady {
			continue
		}
		nftSellsRsp = append(nftSellsRsp, data)
	}

	for _, data := range resp {
		if isReadyOnly && !data.IsReady {
			continue
		}
		nftSellsRsp = append(nftSellsRsp, data)
	}

	return
}

func GetNFTSellUtxoByTokenIndex(key string, tokenIndex string) (nftSellsRsp []*model.NFTSellResp, err error) {
	op := &redis.ZRangeBy{
		Min:    tokenIndex, // 最小分数
		Max:    tokenIndex, // 最大分数
		Offset: 0,          // 类似sql的limit, 表示开始偏移量
		Count:  64,         // 最多兼容64条同样的index
	}
	utxoOutpoints, err := rdb.BizClient.ZRangeByScore(ctx, key, op).Result()
	if err != nil {
		logger.Log.Info("GetUtxoByTokenIndex redis failed", zap.Error(err))
		return
	}
	nftSellsRsp, err = getNFTSellUtxoFromRedis(utxoOutpoints)
	if err != nil {
		return nil, err
	}
	return nftSellsRsp, nil
}
