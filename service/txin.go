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

//////////////// txin
func txInSpentResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxInSpentDO
	err := rows.Scan(&ret.Height, &ret.TxId, &ret.Idx, &ret.UtxId, &ret.Vout)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func txInResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxInDO
	err := rows.Scan(
		&ret.Height, &ret.TxId, &ret.Idx, &ret.ScriptSig,
		&ret.HeightTxo, &ret.UtxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxInputsByTxId(txidHex string) (txInsRsp []*model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
height IN (
    SELECT height FROM tx
    WHERE txid = unhex('%s') LIMIT 1
)`, txidHex, txidHex)

	return GetTxInputsBySql(psql)
}

func GetTxInputsByTxIdInsideHeight(blkHeight int, txidHex string) (txInsRsp []*model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
    height = %d`, txidHex, blkHeight)

	return GetTxInputsBySql(psql)
}

func GetTxInputsByTxIdBeforeHeight(blkHeight int, txidHex string) (txInsRsp []*model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
    height <= %d`, txidHex, blkHeight)

	return GetTxInputsBySql(psql)
}

func GetTxInputsBySql(psql string) (txInsRsp []*model.TxInResp, err error) {
	txInsRet, err := clickhouse.ScanAll(psql, txInResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	txIns := txInsRet.([]*model.TxInDO)
	for _, txin := range txIns {
		txInsRsp = append(txInsRsp, &model.TxInResp{
			Height:       int(txin.Height),
			TxIdHex:      blkparser.HashString(txin.TxId),
			Idx:          int(txin.Idx),
			ScriptSigHex: hex.EncodeToString(txin.ScriptSig),

			HeightTxo:     int(txin.HeightTxo),
			UtxIdHex:      blkparser.HashString(txin.UtxId),
			Vout:          int(txin.Vout),
			Address:       utils.EncodeAddress(txin.Address, utils.PubKeyHashAddrIDMainNet), // fixme
			GenesisHex:    hex.EncodeToString(txin.Genesis),
			Satoshi:       int(txin.Satoshi),
			ScriptTypeHex: hex.EncodeToString(txin.ScriptType),
		})
	}
	return
}

func GetTxInputByTxIdAndIdx(txidHex string, index int) (txInRsp *model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
       idx = %d
LIMIT 1`, txidHex, index)

	return GetTxInputBySql(psql)
}

func GetTxInputByTxIdAndIdxInsideHeight(blkHeight int, txidHex string, index int) (txInRsp *model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
       idx = %d AND
    height = %d
LIMIT 1`, txidHex, index, blkHeight)

	return GetTxInputBySql(psql)
}

func GetTxInputByTxIdAndIdxBeforeHeight(blkHeight int, txidHex string, index int) (txInRsp *model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
       idx = %d AND
    height <= %d
LIMIT 1`, txidHex, index, blkHeight)

	return GetTxInputBySql(psql)
}

func GetTxInputByTxIdAndIdxAfterHeight(blkHeight int, txidHex string, index int) (txInRsp *model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT height, txid, idx, script_sig,
       height_txo, utxid, vout, address, genesis, satoshi, script_type FROM txin_full
WHERE txid = unhex('%s') AND
       idx = %d AND
    height >= %d
LIMIT 1`, txidHex, index, blkHeight)

	return GetTxInputBySql(psql)
}

func GetTxInputBySql(psql string) (txInRsp *model.TxInResp, err error) {
	txInRet, err := clickhouse.ScanOne(psql, txInResultSRF)
	if err != nil {
		log.Printf("query tx by blkid failed: %v", err)
		return nil, err
	}
	if txInRet == nil {
		return nil, errors.New("not exist")
	}

	txin := txInRet.(*model.TxInDO)

	txInRsp = &model.TxInResp{
		Height:       int(txin.Height),
		TxIdHex:      blkparser.HashString(txin.TxId),
		Idx:          int(txin.Idx),
		ScriptSigHex: hex.EncodeToString(txin.ScriptSig),

		HeightTxo:     int(txin.HeightTxo),
		UtxIdHex:      blkparser.HashString(txin.UtxId),
		Vout:          int(txin.Vout),
		Address:       utils.EncodeAddress(txin.Address, utils.PubKeyHashAddrIDMainNet), // fixme
		GenesisHex:    hex.EncodeToString(txin.Genesis),
		Satoshi:       int(txin.Satoshi),
		ScriptTypeHex: hex.EncodeToString(txin.ScriptType),
	}
	return
}

func GetTxOutputSpentStatusByTxIdAndIdx(blkHeight int, txidHex string, index int) (txInRsp *model.TxInSpentResp, err error) {
	var psql string
	if blkHeight < 0 {
		psql = fmt.Sprintf(`
SELECT height, txid, idx, utxid, vout FROM txin_spent
WHERE utxid = unhex('%s') AND
       vout = %d
LIMIT 1`, txidHex, index)
	} else {
		psql = fmt.Sprintf(`
SELECT height, txid, idx, utxid, vout FROM txin_spent
WHERE utxid = unhex('%s') AND
       vout = %d AND
     height = %d
LIMIT 1`,
			txidHex, index, blkHeight)
	}

	txInRet, err := clickhouse.ScanOne(psql, txInSpentResultSRF)
	if err != nil {
		log.Printf("query tx by blkid failed: %v", err)
		return nil, err
	}
	if txInRet == nil {
		return nil, errors.New("not exist")
	}
	txIn := txInRet.(*model.TxInDO)
	txInRsp = &model.TxInSpentResp{
		Height:   int(txIn.Height),
		TxIdHex:  blkparser.HashString(txIn.TxId),
		Idx:      int(txIn.Idx),
		UtxIdHex: blkparser.HashString(txIn.UtxId),
		Vout:     int(txIn.Vout),
	}
	return
}
