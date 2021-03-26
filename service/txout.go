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

const (
	SQL_FIELEDS_TXOUT_WITHOUT_SCRIPT        = "utxid, vout, address, genesis, satoshi, script_type, '', height"
	SQL_FIELEDS_TXOUT_STATUS_WITHOUT_SCRIPT = SQL_FIELEDS_TXOUT_WITHOUT_SCRIPT + ", u.txid, u.height"

	SQL_FIELEDS_TXOUT = "utxid, vout, address, genesis, satoshi, script_type, script_pk, height"
)

//////////////// txout
func txOutStatusResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutStatusDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height,
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
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = unhex('%s')
                    ) AND
           vout >= %d
    ORDER BY vout
    LIMIT %d
) AS u USING (utxid, vout)
WHERE utxid = unhex('%s') AND
     height IN (SELECT height FROM tx_height
                WHERE txid = unhex('%s')
               ) AND
       vout >= %d
ORDER BY vout
LIMIT %d
`, SQL_FIELEDS_TXOUT_STATUS_WITHOUT_SCRIPT,
		txidHex, txidHex, cursor, size,
		txidHex, txidHex, cursor, size)

	return GetTxOutputsBySql(psql)
}

func GetTxOutputsByTxIdInsideHeight(cursor, size, blkHeight int, txidHex string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM txout
LEFT JOIN
(
    SELECT utxid, vout, txid, height FROM txin_spent
    WHERE utxid = unhex('%s') AND
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = unhex('%s')
                    ) AND
           vout >= %d
    ORDER BY vout
    LIMIT %d
) AS u USING (utxid, vout)
WHERE utxid = unhex('%s') AND
     height = %d AND
      vout >= %d
ORDER BY vout
LIMIT %d
`, SQL_FIELEDS_TXOUT_STATUS_WITHOUT_SCRIPT,
		txidHex, txidHex, cursor, size,
		txidHex, blkHeight, cursor, size)

	return GetTxOutputsBySql(psql)
}

func GetTxOutputsBySql(psql string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutStatusResultSRF)
	if err != nil {
		log.Printf("query txouts by blkid failed: %v", err)
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutStatusDO)
	for _, txout := range txOuts {
		address := "-"
		if len(txout.Address) == 20 {
			address = utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet) // fixme
		}
		txOutsRsp = append(txOutsRsp, &model.TxOutStatusResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: address,
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptPkHex:   hex.EncodeToString(txout.ScriptPk),
			Height:        int(txout.Height),

			TxIdSpentHex: blkparser.HashString(txout.TxIdSpent),
			HeightSpent:  int(txout.HeightSpent),
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
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxOutputBySql(psql string) (txOutRsp *model.TxOutResp, err error) {
	txOutRet, err := clickhouse.ScanOne(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txout by blkid failed: %v", err)
		return nil, err
	}
	if txOutRet == nil {
		return nil, errors.New("not exist")
	}
	txOut := txOutRet.(*model.TxOutDO)
	txOutRsp = &model.TxOutResp{
		TxIdHex: blkparser.HashString(txOut.TxId),
		Vout:    int(txOut.Vout),
		Address: utils.EncodeAddress(txOut.Address, utils.PubKeyHashAddrIDMainNet),
		Satoshi: int(txOut.Satoshi),

		GenesisHex:    hex.EncodeToString(txOut.Genesis),
		ScriptTypeHex: hex.EncodeToString(txOut.ScriptType),
		ScriptPkHex:   hex.EncodeToString(txOut.ScriptPk),
		Height:        int(txOut.Height),
	}
	return
}
