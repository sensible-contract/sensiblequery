package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/lib/blkparser"
	"satoblock/lib/utils"
	"satoblock/model"
)

func txOutHistoryResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutHistoryDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.Height, &ret.IOType)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//////////////// address
func GetHistoryByAddress(addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, idx, address, genesis, satoshi, script_type, height, io_type FROM
(
SELECT txid, vout AS idx, address, genesis, satoshi, script_type, height, 1 AS io_type FROM txout
WHERE (txid, vout, height) in (
    SELECT utxid, vout, height FROM txout_address_height
    WHERE address = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)

UNION ALL

SELECT txid, idx, address, genesis, satoshi, script_type, height, 0 AS io_type FROM txin_full
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)
)
ORDER BY height DESC
LIMIT 128
`, addressHex, addressHex)
	return GetHistoryBySql(psql)
}

//////////////// genesis
func GetHistoryByGenesis(genesisHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, idx, address, genesis, satoshi, script_type, height, io_type FROM
(
SELECT txid, vout AS idx, address, genesis, satoshi, script_type, height, 1 AS io_type FROM txout
WHERE (txid, vout, height) in (
    SELECT utxid, vout, height FROM txout_genesis_height
    WHERE genesis = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)

UNION ALL

SELECT txid, idx, address, genesis, satoshi, script_type, height, 0 AS io_type FROM txin_full
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_genesis_height
    WHERE genesis = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)
)
ORDER BY height DESC
LIMIT 128
`, genesisHex, genesisHex)
	return GetHistoryBySql(psql)
}

func GetHistoryBySql(psql string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutHistoryResultSRF)
	if err != nil {
		log.Printf("query txs by genesis failed: %v", err)
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutHistoryDO)
	for _, txout := range txOuts {
		txOutsRsp = append(txOutsRsp, &model.TxOutHistoryResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet), // fixme
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			Height:        int(txout.Height),
			IOType:        int(txout.IOType),
		})
	}
	return
}
