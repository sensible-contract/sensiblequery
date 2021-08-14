package service

import (
	"database/sql"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/logger"
	"sensiblequery/model"

	"go.uber.org/zap"
)

//////////////// history
func txOutHistoryResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutHistoryDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height, &ret.Idx, &ret.IOType)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//////////////// address
func GetHistoryByAddress(addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, io_type FROM
(
SELECT utxid AS txid, vout AS idx, address, codehash, genesis, satoshi, script_type, script_pk, height, utxidx AS txidx, 1 AS io_type FROM txout
WHERE (utxid, vout, height) in (
    SELECT utxid, vout, height FROM txout_address_height
    WHERE address = unhex('%s')
    ORDER BY height DESC
    LIMIT 1024
)

UNION ALL

SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, 0 AS io_type FROM txin
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('%s')
    ORDER BY height DESC
    LIMIT 1024
)
)
ORDER BY height DESC, txidx DESC
LIMIT 1024
`, addressHex, addressHex)
	return GetHistoryBySql(psql)
}

//////////////// genesis
func GetHistoryByGenesis(cursor, size int, codeHashHex, genesisHex, addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	logger.Log.Info("query tx history by codehash/genesis for", zap.String("address", addressHex))
	psql := fmt.Sprintf(`
SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, io_type FROM
(
SELECT utxid AS txid, vout AS idx, address, codehash, genesis, satoshi, script_type, script_pk, height, utxidx AS txidx, 1 AS io_type FROM txout
WHERE (utxid, vout, height) in (
    SELECT utxid, vout, height FROM txout_genesis_height
    WHERE codehash = unhex('%s') AND
          genesis = unhex('%s') AND
          address = unhex('%s')
    ORDER BY height DESC
    LIMIT 1024
)

UNION ALL

SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, 0 AS io_type FROM txin
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_genesis_height
    WHERE codehash = unhex('%s') AND
          genesis = unhex('%s') AND
          address = unhex('%s')
    ORDER BY height DESC
    LIMIT 1024
)
)
ORDER BY height DESC, txidx DESC
LIMIT 1024
`, codeHashHex, genesisHex, addressHex,
		codeHashHex, genesisHex, addressHex)
	return GetHistoryBySql(psql)
}

func GetHistoryBySql(psql string) (txOutHistoriesRsp []*model.TxOutHistoryResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutHistoryResultSRF)
	if err != nil {
		logger.Log.Info("query tx history by genesis failed", zap.Error(err))
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutHistoryDO)
	for _, txout := range txOuts {
		txOutRsp := getTxOutputRespFromDo(&txout.TxOutDO)
		txOutHistoryRsp := &model.TxOutHistoryResp{
			TxOutResp: *txOutRsp,
		}
		txOutHistoryRsp.IOType = int(txout.IOType)
		txOutHistoriesRsp = append(txOutHistoriesRsp, txOutHistoryRsp)
	}
	return
}
