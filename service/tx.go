package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/lib/blkparser"
	"satoblock/model"
)

const (
	SQL_FIELEDS_TX = "txid, nin, nout, txsize, locktime, height, blkid, txidx"
)

//////////////// tx
func txResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxDO
	err := rows.Scan(&ret.TxId, &ret.InCount, &ret.OutCount, &ret.TxSize, &ret.LockTime, &ret.Height, &ret.BlockId, &ret.Idx)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetBlockTxsByBlockHeight(blkHeight int) (txsRsp []*model.TxInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blktx_height WHERE height = %d", SQL_FIELEDS_TX, blkHeight)
	return GetBlockTxsBySql(psql)
}

func GetBlockTxsByBlockId(blkidHex string) (txsRsp []*model.TxInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
WHERE height IN (
    SELECT height FROM blk
    WHERE blkid = unhex('%s') LIMIT 1
)`, SQL_FIELEDS_TX, blkidHex)

	return GetBlockTxsBySql(psql)
}

func GetBlockTxsBySql(psql string) (txsRsp []*model.TxInfoResp, err error) {
	txsRet, err := clickhouse.ScanAll(psql, txResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	if txsRet == nil {
		return nil, errors.New("not exist")
	}
	txs := txsRet.([]*model.TxDO)
	for _, tx := range txs {
		txsRsp = append(txsRsp, &model.TxInfoResp{
			TxIdHex:  blkparser.HashString(tx.TxId),
			InCount:  int(tx.InCount),
			OutCount: int(tx.OutCount),
			TxSize:   int(tx.TxSize),
			LockTime: int(tx.LockTime),
			Height:   int(tx.Height),
			// BlockIdHex: blkparser.HashString(tx.BlockId),
			Idx: int(tx.Idx),
		})
	}
	return
}

func GetTxById(txidHex string) (txRsp *model.TxInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM tx
WHERE txid = unhex('%s') AND
height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)
LIMIT 1`, SQL_FIELEDS_TX, txidHex, txidHex)
	return GetTxBySql(psql)
}

func GetTxByIdInsideHeight(blkHeight int, txidHex string) (txRsp *model.TxInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM tx WHERE txid = unhex('%s') AND height = %d", SQL_FIELEDS_TX, txidHex, blkHeight)
	return GetTxBySql(psql)
}

func GetTxBySql(psql string) (txRsp *model.TxInfoResp, err error) {
	txRet, err := clickhouse.ScanOne(psql, txResultSRF)
	if err != nil {
		log.Printf("query tx failed: %v", err)
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
		Height:     int(tx.Height),
		BlockIdHex: blkparser.HashString(tx.BlockId),
		Idx:        int(tx.Idx),
	}
	return
}
