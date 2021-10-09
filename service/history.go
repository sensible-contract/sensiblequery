package service

import (
	"database/sql"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/logger"
	"sensiblequery/model"

	"go.uber.org/zap"
)

//////////////// history
func txOutHistoryResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutHistoryDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height, &ret.Idx, &ret.IOType, &ret.BlockTime)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//////////////// address
func GetHistoryByAddress(cursor, size int, addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	maxOffset := cursor + size
	psql := fmt.Sprintf(`
SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, io_type, blk.blocktime FROM
(
    SELECT utxid AS txid, vout AS idx, address, codehash, genesis, satoshi, script_type, script_pk, height, utxidx AS txidx, 1 AS io_type FROM txout
    WHERE (utxid, vout, height) in (
        SELECT utxid, vout, height FROM txout_address_height
        WHERE address = unhex('%s')
        ORDER BY height DESC
        LIMIT %d
    )

    UNION ALL

    SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, 0 AS io_type FROM txin
    WHERE (txid, idx, height) in (
        SELECT txid, idx, height FROM txin_address_height
        WHERE address = unhex('%s')
        ORDER BY height DESC
        LIMIT %d
    )
) AS history
LEFT JOIN (
    SELECT height, blocktime FROM blk_height
) AS blk
USING height
ORDER BY height DESC, txidx DESC
LIMIT %d, %d`,
		addressHex, maxOffset,
		addressHex, maxOffset,
		cursor, size)
	return GetHistoryBySql(psql)
}

//////////////// genesis
func GetHistoryByGenesis(cursor, size int, codehashHex, genesisHex, addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	logger.Log.Info("query tx history by codehash/genesis for", zap.String("address", addressHex))
	addressMatch := ""
	if addressHex == "0000000000000000000000000000000000000000" {
		addressMatch = "OR address = ''"
	}

	codehashMatch := ""
	genesisMatch := ""
	if codehashHex != "0000000000000000000000000000000000000000" {
		codehashMatch = fmt.Sprintf("codehash = unhex('%s') AND", codehashHex)
		genesisMatch = fmt.Sprintf("genesis = unhex('%s') AND", genesisHex)
	}

	maxOffset := cursor + size
	// script_pk -> ''
	psql := fmt.Sprintf(`
SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, io_type, blk.blocktime FROM
(
    SELECT utxid AS txid, vout AS idx, address, codehash, genesis, satoshi, script_type, script_pk, height, utxidx AS txidx, 1 AS io_type FROM txout
    WHERE (utxid, vout, height) in (
        SELECT utxid, vout, height FROM txout_genesis_height
        WHERE %s %s (address = unhex('%s') %s)
        ORDER BY height DESC
        LIMIT %d
    )

    UNION ALL

    SELECT txid, idx, address, codehash, genesis, satoshi, script_type, script_pk, height, txidx, 0 AS io_type FROM txin
    WHERE (txid, idx, height) in (
        SELECT txid, idx, height FROM txin_genesis_height
        WHERE %s %s (address = unhex('%s') %s)
        ORDER BY height DESC
        LIMIT %d
    )
) AS history
LEFT JOIN (
    SELECT height, blocktime FROM blk_height
    WHERE height >= 660000
) AS blk
USING height
ORDER BY height DESC, txidx DESC
LIMIT %d, %d`,
		codehashMatch, genesisMatch, addressHex, addressMatch, maxOffset,
		codehashMatch, genesisMatch, addressHex, addressMatch, maxOffset,
		cursor, size)
	return GetHistoryBySql(psql)
}

func GetHistoryBySql(psql string) (txOutHistoriesRsp []*model.TxOutHistoryResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutHistoryResultSRF)
	if err != nil {
		logger.Log.Info("query tx history by genesis failed", zap.Error(err))
		return nil, err
	}
	if txOutsRet == nil {
		// not exist
		txOutHistoriesRsp = make([]*model.TxOutHistoryResp, 0)
		return
	}
	txOuts := txOutsRet.([]*model.TxOutHistoryDO)
	for _, txout := range txOuts {
		txOutRsp := getTxOutputRespFromDo(&txout.TxOutDO)
		txOutRsp.ScriptPkHex = ""
		txOutHistoryRsp := &model.TxOutHistoryResp{
			TxOutResp: *txOutRsp,
		}
		txOutHistoryRsp.IOType = int(txout.IOType)
		txOutHistoryRsp.BlockTime = int(txout.BlockTime)
		txOutHistoriesRsp = append(txOutHistoriesRsp, txOutHistoryRsp)
	}
	return
}
