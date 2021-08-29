package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/lib/blkparser"
	"sensiblequery/lib/utils"
	"sensiblequery/logger"
	"sensiblequery/model"
	"strconv"

	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

const (
	SQL_FIELEDS_TXOUT_WITHOUT_SCRIPT        = "utxid, vout, address, genesis, satoshi, script_type, '', height, txidx"
	SQL_FIELEDS_TXOUT_STATUS_WITHOUT_SCRIPT = SQL_FIELEDS_TXOUT_WITHOUT_SCRIPT + ", u.txid, u.height"

	SQL_FIELEDS_TXOUT        = "utxid, vout, address, codehash, genesis, satoshi, script_type, script_pk, height, utxidx"
	SQL_FIELEDS_TXOUT_STATUS = SQL_FIELEDS_TXOUT + ", u.txid, u.height"
)

//////////////// txout
func txOutStatusResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutStatusDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height, &ret.Idx,
		&ret.TxIdSpent, &ret.HeightSpent)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxOutputsByTxId(cursor, size int, txidHex string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txout
LEFT JOIN
(
    SELECT utxid, vout, txid, height FROM txin_spent
    WHERE utxid = unhex('%s') AND
           vout >= %d AND
        (height = 4294967295 OR
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = unhex('%s') AND
                           vout >= %d
                    ORDER BY vout
                    LIMIT %d
                    ))
    ORDER BY vout
    LIMIT %d
) AS u USING (utxid, vout)
WHERE utxid = unhex('%s') AND
       vout >= %d AND
    (height = 4294967295 OR
     height IN (SELECT height FROM tx_height
                WHERE txid = unhex('%s')
               ))
ORDER BY vout
LIMIT %d
`, SQL_FIELEDS_TXOUT_STATUS, // need without script?
		txidHex, cursor, txidHex, cursor, size, size,
		txidHex, cursor, txidHex, size)

	return GetTxOutputsBySql(psql)
}

func GetTxOutputsByTxIdInsideHeight(cursor, size, blkHeight int, txidHex string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txout
LEFT JOIN
(
    SELECT utxid, vout, txid, height FROM txin_spent
    WHERE utxid = unhex('%s') AND
           vout >= %d AND
         (height == 4294967295 OR
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = unhex('%s') AND
                           vout >= %d
                    ORDER BY vout
                    LIMIT %d
                    ))
    ORDER BY vout
    LIMIT %d
) AS u USING (utxid, vout)
WHERE height = %d AND
     utxid = unhex('%s') AND
      vout >= %d
ORDER BY vout
LIMIT %d
`, SQL_FIELEDS_TXOUT_STATUS, // need without script?
		txidHex, cursor, txidHex, cursor, size, size,
		blkHeight, txidHex, cursor, size)

	return GetTxOutputsBySql(psql)
}

func GetTxOutputsBySql(psql string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutStatusResultSRF)
	if err != nil {
		logger.Log.Info("query txouts by blkid failed", zap.Error(err))
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutStatusDO)
	for _, txout := range txOuts {
		txOutRsp := getTxOutputRespFromDo(&txout.TxOutDO)
		txOutStatusRsp := &model.TxOutStatusResp{
			TxOutResp: *txOutRsp,
		}
		txOutStatusRsp.TxIdSpentHex = blkparser.HashString(txout.TxIdSpent)
		txOutStatusRsp.HeightSpent = int(txout.HeightSpent)

		txOutsRsp = append(txOutsRsp, txOutStatusRsp)
	}
	return
}

func GetTxOutputByTxIdAndIdx(txidHex string, index int) (txOutRsp *model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txout
WHERE utxid = unhex('%s') AND
       vout = %d AND
       (height == 4294967295 OR
        height IN (
            SELECT height FROM tx_height
            WHERE txid = unhex('%s')
       ))
LIMIT 1`, SQL_FIELEDS_TXOUT, txidHex, index, txidHex)
	return GetTxOutputBySql(psql)
}

func GetTxOutputByTxIdAndIdxInsideHeight(blkHeight int, txidHex string, index int) (txOutRsp *model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txout
WHERE utxid = unhex('%s') AND
       vout = %d AND
     height = %d
LIMIT 1`, SQL_FIELEDS_TXOUT, txidHex, index, blkHeight)

	return GetTxOutputBySql(psql)
}

func txOutResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height, &ret.Idx)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxOutputBySql(psql string) (txOutRsp *model.TxOutResp, err error) {
	txOutRet, err := clickhouse.ScanOne(psql, txOutResultSRF)
	if err != nil {
		logger.Log.Info("query txout by blkid failed", zap.Error(err))
		return nil, err
	}
	if txOutRet == nil {
		return nil, errors.New("not exist")
	}
	txout := txOutRet.(*model.TxOutDO)
	txOutRsp = getTxOutputRespFromDo(txout)
	return
}

func getTxOutputRespFromDo(txout *model.TxOutDO) (txOutRsp *model.TxOutResp) {
	txo := scriptDecoder.ExtractPkScriptForTxo(txout.ScriptPk, txout.ScriptType)

	address := ""
	if txo.HasAddress {
		address = utils.EncodeAddress(txo.AddressPkh[:], utils.PubKeyHashAddrID)
	}

	txOutRsp = &model.TxOutResp{
		TxIdHex:       blkparser.HashString(txout.TxId),
		Vout:          int(txout.Vout),
		Address:       address,
		Satoshi:       int(txout.Satoshi),
		ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
		ScriptPkHex:   hex.EncodeToString(txout.ScriptPk),
		Height:        int(txout.Height),
		Idx:           int(txout.Idx),
	}

	txOutRsp.CodeType = int(txo.CodeType)
	if txo.CodeType != scriptDecoder.CodeType_NONE && txo.CodeType != scriptDecoder.CodeType_SENSIBLE {
		txOutRsp.TokenId = hex.EncodeToString(txo.GenesisId[:txo.GenesisIdLen])
		txOutRsp.CodeHashHex = hex.EncodeToString(txo.CodeHash[:])
		txOutRsp.GenesisHex = hex.EncodeToString(txo.GenesisId[:txo.GenesisIdLen])
	}

	if txo.CodeType == scriptDecoder.CodeType_NFT {
		txOutRsp.IsNFT = true
		txOutRsp.TokenIndex = strconv.FormatUint(txo.NFT.TokenIndex, 10)
		txOutRsp.MetaTxIdHex = hex.EncodeToString(txo.NFT.MetaTxId[:])
		txOutRsp.MetaOutputIndex = int(txo.NFT.MetaOutputIndex)
		txOutRsp.SensibleIdHex = hex.EncodeToString(txo.NFT.SensibleId)
	} else if txo.CodeType == scriptDecoder.CodeType_FT {
		txOutRsp.TokenName = txo.FT.Name
		txOutRsp.TokenSymbol = txo.FT.Symbol
		txOutRsp.TokenAmount = strconv.FormatUint(txo.FT.Amount, 10)
		txOutRsp.TokenDecimal = int(txo.FT.Decimal)
		txOutRsp.SensibleIdHex = hex.EncodeToString(txo.FT.SensibleId)
	} else if txo.CodeType == scriptDecoder.CodeType_UNIQUE {
		txOutRsp.SensibleIdHex = hex.EncodeToString(txo.Uniq.SensibleId)
	}
	return
}
