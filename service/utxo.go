package service

import (
	"encoding/hex"
	"sensiblequery/lib/blkparser"
	"sensiblequery/lib/utils"
	"sensiblequery/logger"
	"sensiblequery/model"

	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

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

	// 计算utxo count
	utxoCount, _, _, _, err := GetUtxoCountByAddress(nil, nil, addressPkh, "au")
	if err != nil {
		logger.Log.Info("GetBalanceByAddress utxo count, but redis failed", zap.Error(err))
		return
	}
	balanceRsp.UtxoCount = utxoCount

	return balanceRsp, nil
}

//////////////// address FT utxo count
func GetUtxoCountByAddress(codeHash, genesisId, addressPkh []byte, key string) (total, totalConf, totalUnconf, totalUnconfSpend int, err error) {
	logger.Log.Info("GetUtxoByAddressCount", zap.String("addressHex", hex.EncodeToString(addressPkh)))

	addressKey := ""
	if len(codeHash) == 0 {
		addressKey = string(addressPkh) + "}"
	} else {
		addressKey = addressKey + string(codeHash) + string(genesisId)
	}

	// key: nu, fu, au
	newUtxoKey := "mp:{" + key + addressKey
	addressUtxoConfirmed := "{" + key + addressKey
	addressUtxoSpentUnconfirmed := "mp:s:{" + key + addressKey

	// unconfirmed count
	newUtxoNum, err := rdb.ZCard(ctx, newUtxoKey).Result()
	if err != nil {
		logger.Log.Info("get newUtxoNum from redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("newUtxoNum", zap.Int64("n", newUtxoNum))
	// confirmed count
	addressUtxoConfirmedNum, err := rdb.ZCard(ctx, addressUtxoConfirmed).Result()
	if err != nil {
		logger.Log.Info("get addressUtxoConfirmedNum from redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("addressUtxoConfirmedNum", zap.Int64("n", addressUtxoConfirmedNum))
	// confirmed spending count(spend still unconfirmed)
	addressUtxoSpentUnconfirmedNum, err := rdb.ZCard(ctx, addressUtxoSpentUnconfirmed).Result()
	if err != nil {
		logger.Log.Info("get addressUtxoSpentUnconfirmedNum from redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("addressUtxoSpentUnconfirmedNum", zap.Int64("n", addressUtxoSpentUnconfirmedNum))

	totalConf = int(addressUtxoConfirmedNum)
	totalUnconf = int(newUtxoNum)
	totalUnconfSpend = int(addressUtxoSpentUnconfirmedNum)
	total = totalConf + totalUnconf - totalUnconfSpend

	return total, totalConf, totalUnconf, totalUnconfSpend, nil
}

//////////////// address utxo
func GetUtxoOutpointsByAddress(cursor, size int, codeHash, genesisId, addressPkh []byte, key string) (
	outpoints []string, total, totalConf, totalUnconf, totalUnconfSpend int, err error) {
	logger.Log.Info("GetUtxoOutpointsByAddress", zap.String("addressHex", hex.EncodeToString(addressPkh)))

	addressKey := ""
	if len(codeHash) == 0 {
		addressKey = string(addressPkh) + "}"
	} else {
		addressKey = addressKey + string(codeHash) + string(genesisId)
	}

	// 注意这里查询需要原子化，可使用pipeline
	// key: nu, fu, au
	newUtxoKey := "mp:{" + key + addressKey
	addressUtxoConfirmed := "{" + key + addressKey
	addressUtxoSpentUnconfirmed := "mp:s:{" + key + addressKey
	addressUtxoConfirmedRange := "mp:r:{" + key + addressKey
	tmpUtxoKey := "mp:t:{" + key + addressKey

	total, totalConf, totalUnconf, totalUnconfSpend, err = GetUtxoCountByAddress(codeHash, genesisId, addressPkh, key)
	if err != nil {
		logger.Log.Info("get utxo count from redis failed", zap.Error(err))
		return
	}

	newUtxoOutpoints, err := rdb.ZRevRange(ctx, newUtxoKey, int64(cursor), int64(cursor+size)-1).Result()
	if err == redis.Nil {
		newUtxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetUtxoOutpointsByAddress redis failed", zap.Error(err))
		return
	}
	// 未超过未确认的utxo数量，或者已确认数量减去已花费数量为0，则直接返回
	if cursor+size <= totalUnconf || totalConf == totalUnconfSpend {
		return newUtxoOutpoints, total, totalConf, totalUnconf, totalUnconfSpend, nil
	}

	// 否则需要先提取
	zargs := redis.ZRangeArgs{
		Key:   addressUtxoConfirmed,
		Start: 0,
		Stop:  cursor + size - totalUnconf + totalUnconfSpend,
		Rev:   true,
	}
	nRange, err := rdb.ZRangeStore(ctx, addressUtxoConfirmedRange, zargs).Result()
	if err != nil {
		logger.Log.Info("ZRangeStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZRangeStore", zap.Int64("result", nRange), zap.Any("zargs", zargs))

	// 再去掉已花费的utxo
	nDiff, err := rdb.ZDiffStore(ctx, tmpUtxoKey, addressUtxoConfirmedRange, addressUtxoSpentUnconfirmed).Result()
	if err != nil {
		logger.Log.Info("ZDiffStore redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("ZDiffStore", zap.Int64("n", nDiff))

	// 再提取结果
	utxoOutpoints, err := rdb.ZRevRange(ctx, tmpUtxoKey, int64(cursor), int64(cursor+size-totalUnconf)-1).Result()
	if err == redis.Nil {
		utxoOutpoints = nil
	} else if err != nil {
		logger.Log.Info("GetUtxoOutpointsByAddress redis failed", zap.Error(err))
		return
	}

	for _, utxo := range utxoOutpoints {
		newUtxoOutpoints = append(newUtxoOutpoints, utxo)
	}

	return newUtxoOutpoints, total, totalConf, totalUnconf, totalUnconfSpend, nil
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

//////////////// address utxo
func GetUtxoByAddress(cursor, size int, addressPkh []byte) (
	txOutsRsp []*model.TxStandardOutResp, total, totalConf, totalUnconf, totalUnconfSpend int, err error) {
	logger.Log.Info("GetUtxoByAddress", zap.String("addressHex", hex.EncodeToString(addressPkh)))

	utxoOutpoints, total, totalConf, totalUnconf, totalUnconfSpend, err := GetUtxoOutpointsByAddress(cursor, size, nil, nil, addressPkh, "au")
	if err != nil {
		return
	}

	txOutsRsp, err = getNonTokenUtxoFromRedis(utxoOutpoints)
	return txOutsRsp, total, totalConf, totalUnconf, totalUnconfSpend, err
}
