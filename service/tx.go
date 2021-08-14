package service

import (
	"database/sql"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/lib/blkparser"
	"sensiblequery/logger"
	"sensiblequery/model"

	"go.uber.org/zap"
)

const (
	SQL_FIELEDS_TX           = "txid, nin, nout, txsize, locktime, invalue, outvalue, 0, height, blkid, txidx"
	SQL_FIELEDS_TX_TIMESTAMP = "txid, nin, nout, txsize, locktime, invalue, outvalue, blk.blocktime, height, blkid, txidx"
)

//////////////// tx
func txResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxDO
	err := rows.Scan(&ret.TxId, &ret.InCount, &ret.OutCount, &ret.TxSize, &ret.LockTime, &ret.InSatoshi, &ret.OutSatoshi, &ret.BlockTime, &ret.Height, &ret.BlockId, &ret.Idx)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetBlockTxsByBlockHeight(cursor, size, blkHeight int) (txsRsp []*model.TxInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blktx_height WHERE height = %d AND txidx >= %d ORDER BY txidx LIMIT %d", SQL_FIELEDS_TX, blkHeight, cursor, size)
	return GetBlockTxsBySql(psql)
}

func GetBlockTxsByBlockId(cursor, size int, blkidHex string) (txsRsp []*model.TxInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
WHERE height IN (
    SELECT height FROM blk
    WHERE blkid = unhex('%s') LIMIT 1
) AND txidx >= %d
ORDER BY txidx
LIMIT %d`, SQL_FIELEDS_TX, blkidHex, cursor, size)

	return GetBlockTxsBySql(psql)
}

func GetBlockTxsBySql(psql string) (txsRsp []*model.TxInfoResp, err error) {
	txsRet, err := clickhouse.ScanAll(psql, txResultSRF)
	if err != nil {
		logger.Log.Info("query txs by blkid failed", zap.Error(err))
		return nil, err
	}
	if txsRet == nil {
		return nil, errors.New("not exist")
	}
	txs := txsRet.([]*model.TxDO)
	for _, tx := range txs {
		txsRsp = append(txsRsp, &model.TxInfoResp{
			TxIdHex:    blkparser.HashString(tx.TxId),
			InCount:    int(tx.InCount),
			OutCount:   int(tx.OutCount),
			TxSize:     int(tx.TxSize),
			LockTime:   int(tx.LockTime),
			InSatoshi:  int(tx.InSatoshi),
			OutSatoshi: int(tx.OutSatoshi),

			Height: int(tx.Height),
			// BlockIdHex: blkparser.HashString(tx.BlockId),
			Idx: int(tx.Idx),
		})
	}
	return
}

func GetTxById(txidHex string) (txRsp *model.TxInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
LEFT JOIN  (
    SELECT height, blocktime FROM blk_height
    WHERE height IN (
        SELECT height FROM tx_height
        WHERE txid = unhex('%s')
    )
    LIMIT 1
) AS blk
USING height
WHERE (height = 4294967295 OR
      height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)) AND txid = unhex('%s')
LIMIT 1`, SQL_FIELEDS_TX_TIMESTAMP, txidHex, txidHex, txidHex)
	return GetTxBySql(psql)
}

func GetTxByIdInsideHeight(blkHeight int, txidHex string) (txRsp *model.TxInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
LEFT JOIN (
    SELECT height, blocktime FROM blk_height
    WHERE height = %d
    LIMIT 1
) AS blk
USING height
WHERE height = %d AND txid = unhex('%s')
LIMIT 1`, SQL_FIELEDS_TX_TIMESTAMP, blkHeight, blkHeight, txidHex)
	return GetTxBySql(psql)
}

func GetTxBySql(psql string) (txRsp *model.TxInfoResp, err error) {
	txRet, err := clickhouse.ScanOne(psql, txResultSRF)
	if err != nil {
		logger.Log.Info("query tx failed", zap.Error(err))
		return nil, err
	}
	if txRet == nil {
		return nil, errors.New("not exist")
	}
	tx := txRet.(*model.TxDO)
	txRsp = &model.TxInfoResp{
		TxIdHex:    blkparser.HashString(tx.TxId),
		InCount:    int(tx.InCount),
		OutCount:   int(tx.OutCount),
		TxSize:     int(tx.TxSize),
		LockTime:   int(tx.LockTime),
		InSatoshi:  int(tx.InSatoshi),
		OutSatoshi: int(tx.OutSatoshi),
		BlockTime:  int(tx.BlockTime),
		Height:     int(tx.Height),
		BlockIdHex: blkparser.HashString(tx.BlockId),
		Idx:        int(tx.Idx),
	}
	return
}

////////////////////////////////////////////////////////////////
func rawtxResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret []byte
	err := rows.Scan(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetRawTxById(txidHex string) (txRsp []byte, err error) {
	psql := fmt.Sprintf(`
SELECT rawtx FROM blktx_height
WHERE (height = 4294967295 OR
      height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)) AND txid = unhex('%s')
LIMIT 1`, txidHex, txidHex)

	txRet, err := clickhouse.ScanOne(psql, rawtxResultSRF)
	if err != nil {
		logger.Log.Info("query tx failed", zap.Error(err))
		return nil, err
	}
	if txRet == nil {
		return nil, errors.New("not exist")
	}
	txRsp = txRet.([]byte)

	return txRsp, nil
}

func GetRawTxByIdInsideHeight(blkHeight int, txidHex string) (txRsp []byte, err error) {
	psql := fmt.Sprintf(`
SELECT rawtx FROM blktx_height
WHERE height = %d AND txid = unhex('%s')
LIMIT 1`, blkHeight, txidHex)

	txRet, err := clickhouse.ScanOne(psql, rawtxResultSRF)
	if err != nil {
		logger.Log.Info("query tx failed", zap.Error(err))
		return nil, err
	}
	if txRet == nil {
		return nil, errors.New("not exist")
	}
	txRsp = txRet.([]byte)

	return txRsp, nil
}
