package service

import (
	"encoding/hex"
	"sensiblequery/lib/utils"
	"sensiblequery/logger"
	"sensiblequery/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

////////////////
// ft balance
func GetTokenOwnersByCodeHashGenesis(cursor, size int, codeHash, genesisId []byte) (ftOwnersRsp []*model.FTOwnerBalanceResp, err error) {
	// get decimal from f info
	decimal, err := rdb.HGet(ctx, "fi"+string(codeHash)+string(genesisId), "decimal").Int()
	if err == redis.Nil {
		decimal = 0
	} else if err != nil {
		logger.Log.Info("GetTokenOwnersByCodeHashGenesis decimal, but redis failed", zap.Error(err))
		return
	}

	// merge
	finalKey := "mp:z:{fb" + string(genesisId) + string(codeHash) + "}"

	oldKey := "{fb" + string(genesisId) + string(codeHash) + "}"
	newKey := "mp:{fb" + string(genesisId) + string(codeHash) + "}"

	finalZs := &redis.ZStore{
		Keys: []string{
			oldKey, newKey,
		},
	}

	// 合并已确认余额和未确认余额
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZUnionStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))

	vals, err := rdb.ZRevRangeWithScores(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		logger.Log.Info("GetFTOwnersByCodeHashGenesis redis failed", zap.Error(err))
		return
	}

	pipe := rdb.Pipeline()
	pendingBalanceCmds := make([]*redis.FloatCmd, 0)
	for _, val := range vals {
		logger.Log.Info("GetFTOwnersByCodeHashGenesis", zap.Float64("balance", val.Score))
		pendingBalanceCmds = append(pendingBalanceCmds, pipe.ZScore(ctx, newKey, val.Member.(string)))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for idx, data := range pendingBalanceCmds {
		val := vals[idx]
		balanceRsp := &model.FTOwnerBalanceResp{
			Address: utils.EncodeAddress([]byte(val.Member.(string)), utils.PubKeyHashAddrID),
			Balance: int(val.Score),
			Decimal: decimal,
		}
		ftOwnersRsp = append(ftOwnersRsp, balanceRsp)

		pendingBalance, err := data.Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			panic(err)
		}
		balanceRsp.PendingBalance = int(pendingBalance)
		balanceRsp.Balance -= int(pendingBalance)
	}

	return ftOwnersRsp, nil
}

func GetAllTokenBalanceByAddress(cursor, size int, addressPkh []byte) (ftOwnersRsp []*model.FTSummaryByAddressResp, err error) {
	finalKey := "mp:z:{fs" + string(addressPkh) + "}"

	oldKey := "{fs" + string(addressPkh) + "}"
	newKey := "mp:{fs" + string(addressPkh) + "}"

	finalZs := &redis.ZStore{
		Keys: []string{
			oldKey, newKey,
		},
	}

	// 合并已确认余额和未确认余额
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZUnionStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))

	vals, err := rdb.ZRevRangeWithScores(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		logger.Log.Info("GetAllTokenBalanceByAddress redis failed", zap.Error(err))
		return
	}

	pipe := rdb.Pipeline()
	pendingBalanceCmds := make([]*redis.FloatCmd, 0)
	ftInfoCmds := make([]*redis.StringStringMapCmd, 0)
	for _, val := range vals {
		pendingBalanceCmds = append(pendingBalanceCmds, pipe.ZScore(ctx, newKey, val.Member.(string)))
		// decimal of each token
		ftInfoCmds = append(ftInfoCmds, pipe.HGetAll(ctx, "fi"+val.Member.(string)))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for idx, data := range pendingBalanceCmds {
		val := vals[idx]
		balanceRsp := &model.FTSummaryByAddressResp{
			CodeHashHex: hex.EncodeToString([]byte(val.Member.(string))[:20]),
			GenesisHex:  hex.EncodeToString([]byte(val.Member.(string))[20:]),
			Balance:     int(val.Score),
		}
		ftOwnersRsp = append(ftOwnersRsp, balanceRsp)

		// decimal
		ftinfo, err := ftInfoCmds[idx].Result()
		if err == redis.Nil {
			logger.Log.Info("GetAllTokenBalanceByAddress ftinfo not found")
			ftinfo = map[string]string{
				"decimal":    "0",
				"name":       "",
				"symbol":     "",
				"sensibleid": "",
			}
		} else if err != nil {
			panic(err)
		}

		decimal, _ := strconv.Atoi(ftinfo["decimal"])
		balanceRsp.Decimal = decimal
		balanceRsp.Name = ftinfo["name"]
		balanceRsp.Symbol = ftinfo["symbol"]
		balanceRsp.SensibleIdHex = hex.EncodeToString([]byte(ftinfo["sensibleid"]))

		// balance
		pendingBalance, err := data.Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			panic(err)
		}
		balanceRsp.PendingBalance = int(pendingBalance)
		balanceRsp.Balance -= int(pendingBalance)
	}

	return ftOwnersRsp, nil
}

func GetTokenBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte) (balanceRsp *model.FTOwnerBalanceWithUtxoCountResp, err error) {
	// get decimal from f info
	decimal, err := rdb.HGet(ctx, "fi"+string(codeHash)+string(genesisId), "decimal").Int()
	if err == redis.Nil {
		decimal = 0
	} else if err != nil {
		logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress decimal, but redis failed", zap.Error(err))
		return
	}

	balance, err := rdb.ZScore(ctx, "{fb"+string(genesisId)+string(codeHash)+"}", string(addressPkh)).Result()
	if err == redis.Nil {
		logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress fb, but not found")
		balance = 0
	} else if err != nil {
		logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress fb, but redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress fb", zap.Float64("balance", balance))
	mpBalance, err := rdb.ZScore(ctx, "mp:{fb"+string(genesisId)+string(codeHash)+"}", string(addressPkh)).Result()
	if err == redis.Nil {
		logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress mp:fb, but not found")
		mpBalance = 0
	} else if err != nil {
		logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress mp:fb, but redis mp failed", zap.Error(err))
		return
	}

	// 计算utxo count，有点费
	logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress fb", zap.Float64("pendingBalance", mpBalance))
	finalUtxoKey, err := mergeUtxoByCodeHashGenesisAddress(codeHash, genesisId, addressPkh, false)
	if err != nil {
		return
	}

	utxoCount, err := rdb.ZCard(ctx, finalUtxoKey).Result()
	if err != nil {
		logger.Log.Info("GetTokenBalanceByCodeHashGenesisAddress merge, but redis failed", zap.Error(err))
		return
	}

	balanceRsp = &model.FTOwnerBalanceWithUtxoCountResp{
		Address:        utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrID),
		Balance:        int(balance),
		PendingBalance: int(mpBalance),
		UtxoCount:      int(utxoCount),
		Decimal:        decimal,
	}
	return balanceRsp, nil
}

////////////////
// nft
func GetNFTOwnersByCodeHashGenesis(cursor, size int, codeHash, genesisId []byte) (ownersRsp []*model.NFTOwnerResp, err error) {
	finalKey := "mp:z:{no" + string(genesisId) + string(codeHash) + "}"

	oldKey := "{no" + string(genesisId) + string(codeHash) + "}"
	newKey := "mp:{no" + string(genesisId) + string(codeHash) + "}"

	finalZs := &redis.ZStore{
		Keys: []string{
			oldKey, newKey,
		},
	}

	// 合并已确认数量和未确认数量
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZUnionStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))

	vals, err := rdb.ZRevRangeWithScores(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		logger.Log.Info("GetNFTOwnersByCodeHashGenesis redis failed", zap.Error(err))
		return
	}

	pipe := rdb.Pipeline()
	pendingCountCmds := make([]*redis.FloatCmd, 0)
	for _, val := range vals {
		pendingCountCmds = append(pendingCountCmds, pipe.ZScore(ctx, newKey, val.Member.(string)))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for idx, data := range pendingCountCmds {
		val := vals[idx]

		countRsp := &model.NFTOwnerResp{
			Address: utils.EncodeAddress([]byte(val.Member.(string)), utils.PubKeyHashAddrID),
			Count:   int(val.Score),
		}
		ownersRsp = append(ownersRsp, countRsp)

		pendingCount, err := data.Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			panic(err)
		}
		countRsp.PendingCount = int(pendingCount)
		countRsp.Count -= int(pendingCount)

	}

	return ownersRsp, nil
}

func GetAllNFTBalanceByAddress(cursor, size int, addressPkh []byte) (nftOwnersRsp []*model.NFTSummaryByAddressResp, err error) {
	finalKey := "mp:z:{ns" + string(addressPkh) + "}"

	oldKey := "{ns" + string(addressPkh) + "}"
	newKey := "mp:{ns" + string(addressPkh) + "}"

	finalZs := &redis.ZStore{
		Keys: []string{
			oldKey, newKey,
		},
	}

	// 合并已确认数量和未确认数量
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		logger.Log.Info("ZUnionStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZUnionStore", zap.Int64("n", nUnion))

	vals, err := rdb.ZRevRangeWithScores(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		logger.Log.Info("GetAllNFTBalanceByAddress redis failed", zap.Error(err))
		return
	}

	pipe := rdb.Pipeline()
	pendingCountCmds := make([]*redis.FloatCmd, 0)
	nftInfoCmds := make([]*redis.StringStringMapCmd, 0)
	for _, val := range vals {
		pendingCountCmds = append(pendingCountCmds, pipe.ZScore(ctx, newKey, val.Member.(string)))
		// metatx of each token
		nftInfoCmds = append(nftInfoCmds, pipe.HGetAll(ctx, "nI"+val.Member.(string)+"1"))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for idx, data := range pendingCountCmds {
		val := vals[idx]
		countRsp := &model.NFTSummaryByAddressResp{
			CodeHashHex: hex.EncodeToString([]byte(val.Member.(string))[:20]),
			GenesisHex:  hex.EncodeToString([]byte(val.Member.(string))[20:]),
			Count:       int(val.Score),
		}
		nftOwnersRsp = append(nftOwnersRsp, countRsp)

		// sensible/supply
		if nftinfo, err := nftInfoCmds[idx].Result(); err == nil {
			supply, _ := strconv.Atoi(nftinfo["supply"])
			countRsp.Supply = supply
			metavout, _ := strconv.Atoi(nftinfo["metavout"])
			countRsp.MetaOutputIndex = metavout
			countRsp.MetaTxIdHex = hex.EncodeToString([]byte(nftinfo["metatxid"]))
			countRsp.SensibleIdHex = hex.EncodeToString([]byte(nftinfo["sensibleid"]))
		} else if err == redis.Nil {
			logger.Log.Info("GetAllTokenBalanceByAddress ftinfo not found")
		} else if err != nil {
			panic(err)
		}

		// count
		pendingCount, err := data.Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			panic(err)
		}
		countRsp.PendingCount = int(pendingCount)
		countRsp.Count -= int(pendingCount)
	}

	return nftOwnersRsp, nil
}

func GetNFTCountByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte) (countRsp *model.NFTOwnerResp, err error) {
	score, err := rdb.ZScore(ctx, "{no"+string(genesisId)+string(codeHash)+"}", string(addressPkh)).Result()
	if err == redis.Nil {
		score = 0
	} else if err != nil {
		logger.Log.Info("GetNFTCountByCodeHashGenesisAddress redis failed", zap.Error(err))
		return
	}

	mpScore, err := rdb.ZScore(ctx, "mp:{no"+string(genesisId)+string(codeHash)+"}", string(addressPkh)).Result()
	if err == redis.Nil {
		mpScore = 0
	} else if err != nil {
		logger.Log.Info("GetNFTCountByCodeHashGenesisAddress mp redis failed", zap.Error(err))
		return
	}

	countRsp = &model.NFTOwnerResp{
		Address:      utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrID),
		Count:        int(score),
		PendingCount: int(mpScore),
	}
	return countRsp, nil
}
