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

//////////////// tx
func txResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxDO
	err := rows.Scan(&ret.TxId, &ret.InCount, &ret.OutCount, &ret.Height, &ret.BlockId, &ret.Idx)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetBlockTxsByBlockHeight(blkHeight int) (txsRsp []*model.TxInfoResp, err error) {
	psql := fmt.Sprintf("SELECT txid, nin, nout, height, blkid, idx FROM blktx WHERE height = %d", blkHeight)
	return GetBlockTxsBySql(psql)
}

func GetBlockTxsByBlockId(blkidHex string) (txsRsp []*model.TxInfoResp, err error) {
	psql := fmt.Sprintf("SELECT txid, nin, nout, height, blkid, idx FROM blktx WHERE blkid = unhex('%s')", blkidHex)
	return GetBlockTxsBySql(psql)
}

func GetBlockTxsBySql(psql string) (txsRsp []*model.TxInfoResp, err error) {
	txsRet, err := clickhouse.ScanAll(psql, txResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	txs := txsRet.([]*model.TxDO)
	for _, tx := range txs {
		txsRsp = append(txsRsp, &model.TxInfoResp{
			TxIdHex:  blkparser.HashString(tx.TxId),
			InCount:  int(tx.InCount),
			OutCount: int(tx.OutCount),

			Height: int(tx.Height),
			// BlockIdHex: blkparser.HashString(tx.BlockId),
			Idx: int(tx.Idx),
		})
	}
	return
}

func GetTxById(blkHeight int, txidHex string) (txRsp *model.TxInfoResp, err error) {
	var psql string
	if blkHeight < 0 {
		psql = fmt.Sprintf("SELECT txid, nin, nout, height, blkid, idx FROM tx WHERE txid = unhex('%s')", txidHex)
	} else {
		psql = fmt.Sprintf("SELECT txid, nin, nout, height, blkid, idx FROM tx WHERE txid = unhex('%s') AND height = %d", txidHex, blkHeight)
	}

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
		TxIdHex:  blkparser.HashString(tx.TxId),
		InCount:  int(tx.InCount),
		OutCount: int(tx.OutCount),

		Height:     int(tx.Height),
		BlockIdHex: blkparser.HashString(tx.BlockId),
		Idx:        int(tx.Idx),
	}
	return
}
