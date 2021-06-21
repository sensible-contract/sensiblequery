package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"satosensible/dao/clickhouse"
	"satosensible/lib/blkparser"
	"satosensible/lib/script"
	"satosensible/lib/utils"
	"satosensible/model"
	"strconv"
)

const (
	SQL_FIELEDS_TXIN_WITHOUT_SCRIPT = `height, txid, idx, '', nsequence,
       height_txo, utxid, vout, address, genesis, satoshi, script_type, ''`

	SQL_FIELEDS_TXIN = `height, txid, idx, script_sig, nsequence,
       height_txo, utxid, vout, address, codehash, genesis, satoshi, script_type, script_pk`
	SQL_FIELEDS_TXIN_SPENT = "height, txid, idx, utxid, vout"
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
		&ret.Height, &ret.TxId, &ret.Idx, &ret.ScriptSig, &ret.Sequence,
		&ret.HeightTxo, &ret.UtxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxInputsByTxId(cursor, size int, txidHex string) (txInsRsp []*model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txin
WHERE txid = unhex('%s') AND
   (height = 4294967295 OR
    height IN (
        SELECT height FROM tx_height
        WHERE txid = unhex('%s')
    )) AND
      idx >= %d
ORDER BY idx
LIMIT %d`, SQL_FIELEDS_TXIN, txidHex, txidHex, cursor, size)

	return GetTxInputsBySql(psql)
}

func GetTxInputsByTxIdInsideHeight(cursor, size, blkHeight int, txidHex string) (txInsRsp []*model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txin
WHERE txid = unhex('%s') AND
    height = %d AND
      idx >= %d
ORDER BY idx
LIMIT %d`, SQL_FIELEDS_TXIN, txidHex, blkHeight, cursor, size)

	return GetTxInputsBySql(psql)
}

func GetTxInputsBySql(psql string) (txInsRsp []*model.TxInResp, err error) {
	txInsRet, err := clickhouse.ScanAll(psql, txInResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	if txInsRet == nil {
		return nil, errors.New("not exist")
	}
	txIns := txInsRet.([]*model.TxInDO)
	for _, txin := range txIns {
		address := "-"
		if len(txin.Address) == 20 {
			address = utils.EncodeAddress(txin.Address, utils.PubKeyHashAddrID)
		}
		isNFT, _, _, _, metaTxId, name, symbol, dataValue, decimal := script.ExtractPkScriptForTxo(txin.ScriptPk, txin.ScriptType)
		tokenId := ""
		if len(txin.Genesis) >= 20 {
			if isNFT {
				tokenId = strconv.Itoa(int(dataValue))
			} else {
				tokenId = hex.EncodeToString(txin.Genesis)
			}
		}
		txInsRsp = append(txInsRsp, &model.TxInResp{
			Height:       int(txin.Height),
			TxIdHex:      blkparser.HashString(txin.TxId),
			Idx:          int(txin.Idx),
			ScriptSigHex: hex.EncodeToString(txin.ScriptSig),
			Sequence:     int(txin.Sequence),

			HeightTxo:     int(txin.HeightTxo),
			UtxIdHex:      blkparser.HashString(txin.UtxId),
			Vout:          int(txin.Vout),
			Address:       address,
			IsNFT:         isNFT,
			TokenId:       tokenId,
			MetaTxId:      hex.EncodeToString(metaTxId),
			TokenName:     name,
			TokenSymbol:   symbol,
			TokenAmount:   strconv.Itoa(int(dataValue)),
			TokenDecimal:  int(decimal),
			CodeHashHex:   hex.EncodeToString(txin.CodeHash),
			GenesisHex:    hex.EncodeToString(txin.Genesis),
			Satoshi:       int(txin.Satoshi),
			ScriptTypeHex: hex.EncodeToString(txin.ScriptType),
			ScriptPkHex:   hex.EncodeToString(txin.ScriptPk),
		})
	}
	return
}

func GetTxInputByTxIdAndIdx(txidHex string, index int) (txInRsp *model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txin
WHERE txid = unhex('%s') AND
       idx = %d AND
height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)
LIMIT 1`, SQL_FIELEDS_TXIN, txidHex, index, txidHex)

	return GetTxInputBySql(psql)
}

func GetTxInputByTxIdAndIdxInsideHeight(blkHeight int, txidHex string, index int) (txInRsp *model.TxInResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txin
WHERE txid = unhex('%s') AND
       idx = %d AND
    height = %d
LIMIT 1`, SQL_FIELEDS_TXIN, txidHex, index, blkHeight)

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

	address := "-"
	if len(txin.Address) == 20 {
		address = utils.EncodeAddress(txin.Address, utils.PubKeyHashAddrID)
	}
	isNFT, _, _, _, metaTxId, name, symbol, dataValue, decimal := script.ExtractPkScriptForTxo(txin.ScriptPk, txin.ScriptType)
	tokenId := ""
	if len(txin.Genesis) >= 20 {
		if isNFT {
			tokenId = strconv.Itoa(int(dataValue))
		} else {
			tokenId = hex.EncodeToString(txin.Genesis)
		}
	}

	txInRsp = &model.TxInResp{
		Height:       int(txin.Height),
		TxIdHex:      blkparser.HashString(txin.TxId),
		Idx:          int(txin.Idx),
		ScriptSigHex: hex.EncodeToString(txin.ScriptSig),

		HeightTxo:     int(txin.HeightTxo),
		UtxIdHex:      blkparser.HashString(txin.UtxId),
		Vout:          int(txin.Vout),
		Address:       address,
		IsNFT:         isNFT,
		TokenId:       tokenId,
		MetaTxId:      hex.EncodeToString(metaTxId),
		TokenName:     name,
		TokenSymbol:   symbol,
		TokenAmount:   strconv.Itoa(int(dataValue)),
		TokenDecimal:  int(decimal),
		CodeHashHex:   hex.EncodeToString(txin.CodeHash),
		GenesisHex:    hex.EncodeToString(txin.Genesis),
		Satoshi:       int(txin.Satoshi),
		ScriptTypeHex: hex.EncodeToString(txin.ScriptType),
	}
	return
}

func GetTxOutputSpentStatusByTxIdAndIdx(txidHex string, index int) (txInRsp *model.TxInSpentResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txin_spent
WHERE utxid = unhex('%s') AND
       vout = %d AND
height IN (
    SELECT height FROM txout_spent_height
    WHERE utxid = unhex('%s') AND
           vout = %d
)
LIMIT 1`, SQL_FIELEDS_TXIN_SPENT, txidHex, index, txidHex, index)

	txInRet, err := clickhouse.ScanOne(psql, txInSpentResultSRF)
	if err != nil {
		log.Printf("query tx by blkid failed: %v", err)
		return nil, err
	}
	if txInRet == nil {
		return nil, errors.New("not exist")
	}
	txIn := txInRet.(*model.TxInSpentDO)
	txInRsp = &model.TxInSpentResp{
		Height:   int(txIn.Height),
		TxIdHex:  blkparser.HashString(txIn.TxId),
		Idx:      int(txIn.Idx),
		UtxIdHex: blkparser.HashString(txIn.UtxId),
		Vout:     int(txIn.Vout),
	}
	return
}
