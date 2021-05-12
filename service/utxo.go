package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"satosensible/lib/blkparser"
	"satosensible/lib/script"
	"satosensible/lib/utils"
	"satosensible/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
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
		return balanceRsp, nil
	} else if err != nil {
		log.Printf("GetBalanceByAddress redis failed: %v", err)
		return
	}
	balanceRsp.Satoshi = int(balance)
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
		log.Printf("GetUtxoByTokenId redis failed: %v", err)
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

//////////////// genesisId
func GetUtxoByCodeHashGenesisAddress(cursor, size int, codeHash, genesisId, addressPkh []byte, isNFT bool) (txOutsRsp []*model.TxOutResp, err error) {
	// 注意这里查询需要原子化，可使用pipeline
	addressKey := string(codeHash) + string(genesisId) + string(addressPkh)
	addressUtxoConfirmed := ""
	addressUtxoSpentUnconfirmed := ""
	oldUtxoKey := ""
	newUtxoKey := ""
	finalKey := ""

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
		log.Printf("ZDiffStore redis failed: %v", err)
		return
	}
	log.Printf("ZDiffStore : %v", nDiff)

	finalZs := &redis.ZStore{
		Keys: []string{
			oldUtxoKey, newUtxoKey,
		},
	}
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		log.Printf("ZDiffStore redis failed: %v", err)
		return
	}
	log.Printf("ZUnionStore : %v", nUnion)

	//
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetUtxoByCodeHashGenesisAddress redis failed: %v", err)
		return
	}
	return getUtxoFromRedis(utxoOutpoints)
}

////////////////
func getUtxoFromRedis(utxoOutpoints []string) (txOutsRsp []*model.TxOutResp, err error) {
	log.Printf("getUtxoFromRedis redis: %d", len(utxoOutpoints))
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
			log.Printf("redis not found outpoint: %s", hex.EncodeToString([]byte(outpoint)))
			continue
		} else if err != nil {
			panic(err)
		}
		txout := &model.TxoData{}
		txout.Unmarshal([]byte(res))

		// 补充数据
		txout.UTxid = []byte(outpoint[:32])                            // 32
		txout.Vout = binary.LittleEndian.Uint32([]byte(outpoint[32:])) // 4
		txout.ScriptType = script.GetLockingScriptType(txout.Script)
		txout.IsNFT, txout.CodeHash, txout.GenesisId, txout.AddressPkh, txout.DataValue, txout.Decimal = script.ExtractPkScriptForTxo(txout.Script, txout.ScriptType)

		tokenId := ""
		if len(txout.GenesisId) >= 20 {
			if txout.IsNFT {
				tokenId = strconv.Itoa(int(txout.DataValue))
			} else {
				tokenId = hex.EncodeToString(txout.GenesisId)
			}
		}

		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.AddressPkh, utils.PubKeyHashAddrID),
			Satoshi: int(txout.Satoshi),

			IsNFT:         txout.IsNFT,
			TokenId:       tokenId,
			TokenAmount:   int(txout.DataValue),
			TokenDecimal:  int(txout.Decimal),
			CodeHashHex:   hex.EncodeToString(txout.CodeHash),
			GenesisHex:    hex.EncodeToString(txout.GenesisId),
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
	log.Printf("GetUtxoByAddress: %s", hex.EncodeToString(addressPkh))

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
		log.Printf("ZDiffStore redis failed: %v", err)
		return
	}
	log.Printf("ZDiffStore : %v", nDiff)

	finalZs := &redis.ZStore{
		Keys: []string{
			oldUtxoKey, newUtxoKey,
		},
	}
	nUnion, err := rdb.ZUnionStore(ctx, finalKey, finalZs).Result()
	if err != nil {
		log.Printf("ZDiffStore redis failed: %v", err)
		return
	}
	log.Printf("ZUnionStore : %v", nUnion)
	utxoOutpoints, err := rdb.ZRevRange(ctx, finalKey, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetUtxoByAddress redis failed: %v", err)
		return
	}
	return getNonTokenUtxoFromRedis(utxoOutpoints)
}

////////////////
func getNonTokenUtxoFromRedis(utxoOutpoints []string) (txOutsRsp []*model.TxStandardOutResp, err error) {
	log.Printf("getUtxoFromRedis redis: %d", len(utxoOutpoints))
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
			log.Printf("redis not found outpoint: %s", hex.EncodeToString([]byte(outpoint)))
			continue
		} else if err != nil {
			panic(err)
		}
		txout := &model.TxoData{}
		txout.Unmarshal([]byte(res))

		// 补充数据
		txout.UTxid = []byte(outpoint[:32])                            // 32
		txout.Vout = binary.LittleEndian.Uint32([]byte(outpoint[32:])) // 4
		txout.ScriptType = script.GetLockingScriptType(txout.Script)

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
