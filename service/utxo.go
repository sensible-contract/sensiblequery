package service

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/lib/blkparser"
	"satoblock/lib/script"
	"satoblock/lib/utils"
	"satoblock/model"

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

//////////////// address
func GetUtxoByAddress(addressPkh []byte) (txOutsRsp []*model.TxOutResp, err error) {
	vals, err := rdb.ZRange("a"+string(addressPkh), 0, 64).Result()
	if err != nil {
		panic(err)
	}

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
		txout.GenesisId, txout.AddressPkh = script.ExtractPkScriptAddressPkh(txout.Script, txout.ScriptType)
		if txout.AddressPkh == nil {
			txout.GenesisId, txout.AddressPkh = script.ExtractPkScriptGenesisIdAndAddressPkh(txout.Script)
		}

		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.AddressPkh, utils.PubKeyHashAddrIDMainNet),
			Satoshi: int(txout.Value),

			GenesisHex:    hex.EncodeToString(txout.GenesisId),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptPkHex:   hex.EncodeToString(txout.Script),
			Height:        int(txout.BlockHeight),
			Idx:           int(txout.TxIdx),
		})
	}

	return txOutsRsp, nil
}

//////////////// genesis
func GetUtxoByGenesis(genesisHex string) (txOutsRsp []*model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM utxo_genesis
WHERE genesis = unhex('%s')
ORDER BY height DESC
LIMIT 128
`, SQL_FIELEDS_TXOUT, genesisHex)
	return GetUtxoBySql(psql)
}

func GetUtxoBySql(psql string) (txOutsRsp []*model.TxOutResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query utxo by sql failed: %v", err)
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutDO)
	for _, txout := range txOuts {
		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet),
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptPkHex:   hex.EncodeToString(txout.ScriptPk),
			Height:        int(txout.Height),
		})
	}
	return
}
