package service

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"satosensible/lib/blkparser"
	"satosensible/lib/script"
	"satosensible/lib/utils"
	"satosensible/model"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var rdb *redis.Client

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
}

func GetBalanceByAddress(addressPkh []byte) (balanceRsp *model.BalanceResp, err error) {
	balanceRsp = &model.BalanceResp{
		Address: utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrID),
	}

	balance, err := rdb.ZScore("balance", string(addressPkh)).Result()
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
func GetUtxoByAddress(cursor, size int, addressPkh []byte) (txOutsRsp []*model.TxOutResp, err error) {
	return getUtxoFromRedis(cursor, size, "au"+string(addressPkh))
}

//////////////// genesisId
func GetUtxoByCodeHashGenesisAddress(cursor, size int, codeHash, genesisId, addressPkh []byte, isNFT bool) (txOutsRsp []*model.TxOutResp, err error) {
	if isNFT {
		return getUtxoFromRedis(cursor, size, "nu"+string(codeHash)+string(genesisId)+string(addressPkh))
	} else {
		return getUtxoFromRedis(cursor, size, "fu"+string(codeHash)+string(genesisId)+string(addressPkh))
	}
}

////////////////
func getUtxoFromRedis(cursor, size int, key string) (txOutsRsp []*model.TxOutResp, err error) {
	vals, err := rdb.ZRevRange(key, int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetUtxoByAddress redis failed: %v", err)
		return
	}

	log.Printf("getUtxoFromRedis redis: %d", len(vals))
	pipe := rdb.Pipeline()
	m := map[string]*redis.StringCmd{}
	for _, key := range vals {
		m[key] = pipe.Get(key)
	}
	_, err = pipe.Exec()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for key, v := range m {
		res, err := v.Result()
		if err == redis.Nil {
			log.Printf("redis not found key: %s", hex.EncodeToString([]byte(key)))
			continue
		} else if err != nil {
			panic(err)
		}
		txout := &model.TxoData{}
		txout.Unmarshal([]byte(res))

		// 补充数据
		txout.UTxid = []byte(key[:32])                            // 32
		txout.Vout = binary.LittleEndian.Uint32([]byte(key[32:])) // 4
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
			ScriptPkHex:   hex.EncodeToString(txout.Script),
			Height:        int(txout.BlockHeight),
			Idx:           int(txout.TxIdx),
		})
	}

	return txOutsRsp, nil
}
