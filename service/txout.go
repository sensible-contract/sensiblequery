package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"satosensible/dao/clickhouse"
	"satosensible/lib/blkparser"
	"satosensible/lib/utils"
	"satosensible/logger"
	"satosensible/model"
	"strconv"

	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

const (
	SQL_FIELEDS_TXOUT_WITHOUT_SCRIPT        = "utxid, vout, address, genesis, satoshi, script_type, '', height"
	SQL_FIELEDS_TXOUT_STATUS_WITHOUT_SCRIPT = SQL_FIELEDS_TXOUT_WITHOUT_SCRIPT + ", u.txid, u.height"

	SQL_FIELEDS_TXOUT        = "utxid, vout, address, codehash, genesis, satoshi, script_type, script_pk, height"
	SQL_FIELEDS_TXOUT_STATUS = SQL_FIELEDS_TXOUT + ", u.txid, u.height"
)

//////////////// txout
func txOutStatusResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutStatusDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height,
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
	for _, txOut := range txOuts {
		address := "-"
		if len(txOut.Address) == 20 {
			address = utils.EncodeAddress(txOut.Address, utils.PubKeyHashAddrID)
		}

		txo := scriptDecoder.ExtractPkScriptForTxo(txOut.ScriptPk, txOut.ScriptType)

		txOutsRsp = append(txOutsRsp, &model.TxOutStatusResp{
			TxIdHex: blkparser.HashString(txOut.TxId),
			Vout:    int(txOut.Vout),
			Address: address,
			Satoshi: int(txOut.Satoshi),

			IsNFT:           (txo.CodeType == scriptDecoder.CodeType_NFT),
			CodeType:        int(txo.CodeType),
			TokenIndex:      strconv.FormatUint(txo.TokenIndex, 10),
			MetaTxIdHex:     hex.EncodeToString(txo.MetaTxId),
			MetaOutputIndex: int(txo.MetaOutputIndex),
			TokenId:         hex.EncodeToString(txOut.Genesis),
			TokenName:       txo.Name,
			TokenSymbol:     txo.Symbol,
			TokenAmount:     strconv.FormatUint(txo.Amount, 10),
			TokenDecimal:    int(txo.Decimal),
			CodeHashHex:     hex.EncodeToString(txOut.CodeHash),
			GenesisHex:      hex.EncodeToString(txOut.Genesis),
			SensibleIdHex:   hex.EncodeToString(txo.SensibleId),
			ScriptTypeHex:   hex.EncodeToString(txOut.ScriptType),
			ScriptPkHex:     hex.EncodeToString(txOut.ScriptPk),
			Height:          int(txOut.Height),

			TxIdSpentHex: blkparser.HashString(txOut.TxIdSpent),
			HeightSpent:  int(txOut.HeightSpent),
		})
	}
	return
}

func GetTxOutputByTxIdAndIdx(txidHex string, index int) (txOutRsp *model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txout
WHERE utxid = unhex('%s') AND
       vout = %d AND
height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)
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
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.CodeHash, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height)
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
	txOut := txOutRet.(*model.TxOutDO)
	address := "-"
	if len(txOut.Address) == 20 {
		address = utils.EncodeAddress(txOut.Address, utils.PubKeyHashAddrID)
	}

	txo := scriptDecoder.ExtractPkScriptForTxo(txOut.ScriptPk, txOut.ScriptType)

	txOutRsp = &model.TxOutResp{
		TxIdHex: blkparser.HashString(txOut.TxId),
		Vout:    int(txOut.Vout),
		Address: address,
		Satoshi: int(txOut.Satoshi),

		IsNFT:           (txo.CodeType == scriptDecoder.CodeType_NFT),
		CodeType:        int(txo.CodeType),
		TokenIndex:      strconv.FormatUint(txo.TokenIndex, 10),
		MetaTxIdHex:     hex.EncodeToString(txo.MetaTxId),
		MetaOutputIndex: int(txo.MetaOutputIndex),
		TokenId:         hex.EncodeToString(txOut.Genesis),
		TokenName:       txo.Name,
		TokenSymbol:     txo.Symbol,
		TokenAmount:     strconv.FormatUint(txo.Amount, 10),
		TokenDecimal:    int(txo.Decimal),
		CodeHashHex:     hex.EncodeToString(txOut.CodeHash),
		GenesisHex:      hex.EncodeToString(txOut.Genesis),
		ScriptTypeHex:   hex.EncodeToString(txOut.ScriptType),
		ScriptPkHex:     hex.EncodeToString(txOut.ScriptPk),
		Height:          int(txOut.Height),
	}
	return
}
