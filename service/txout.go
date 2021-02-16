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

//////////////// txout
func txOutStatusResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutStatusDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.Script, &ret.Height,
		&ret.TxIdSpent, &ret.HeightSpent)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxOutputsByTxId(txidHex string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height, u.txid, u.height FROM txout
LEFT JOIN
(
    SELECT utxid, vout, txid, height FROM txin_spent
    WHERE utxid = unhex('%s') AND
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = unhex('%s')
                    )
) AS u ON txout.txid = u.utxid AND txout.vout = u.vout
WHERE txid = unhex('%s') AND
height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)
`, txidHex, txidHex, txidHex, txidHex)

	return GetTxOutputsBySql(psql)
}

func GetTxOutputsByTxIdInsideHeight(blkHeight int, txidHex string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height, u.txid, u.height FROM txout
LEFT JOIN
(
    SELECT utxid, vout, txid, height FROM txin_spent
    WHERE utxid = unhex('%s') AND
         height IN (SELECT height FROM txout_spent_height
                    WHERE utxid = unhex('%s')
                    )
) AS u ON txout.txid = u.utxid AND txout.vout = u.vout
WHERE txid = unhex('%s') AND
    height = %d
`, txidHex, txidHex, txidHex, blkHeight)

	return GetTxOutputsBySql(psql)
}

func GetTxOutputsBySql(psql string) (txOutsRsp []*model.TxOutStatusResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutStatusResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutStatusDO)
	for _, txout := range txOuts {
		txOutsRsp = append(txOutsRsp, &model.TxOutStatusResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet), // fixme
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptHex:     hex.EncodeToString(txout.Script),
			Height:        int(txout.Height),

			TxIdSpentHex: blkparser.HashString(txout.TxIdSpent),
			HeightSpent:  int(txout.HeightSpent),
		})
	}
	return
}

func GetTxOutputByTxIdAndIdx(txidHex string, index int) (txOutRsp *model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout
WHERE txid = unhex('%s') AND
      vout = %d AND
height IN (
    SELECT height FROM tx_height
    WHERE txid = unhex('%s')
)
LIMIT 1`, txidHex, index, txidHex)
	return GetTxOutputBySql(psql)
}

func GetTxOutputByTxIdAndIdxInsideHeight(blkHeight int, txidHex string, index int) (txOutRsp *model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout
WHERE txid = unhex('%s') AND
      vout = %d AND
    height = %d
LIMIT 1`, txidHex, index, blkHeight)

	return GetTxOutputBySql(psql)
}

func txOutResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.Script, &ret.Height)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxOutputBySql(psql string) (txOutRsp *model.TxOutResp, err error) {
	txOutRet, err := clickhouse.ScanOne(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	if txOutRet == nil {
		return nil, errors.New("not exist")
	}
	txOut := txOutRet.(*model.TxOutDO)
	txOutRsp = &model.TxOutResp{
		TxIdHex: blkparser.HashString(txOut.TxId),
		Vout:    int(txOut.Vout),
		Address: utils.EncodeAddress(txOut.Address, utils.PubKeyHashAddrIDMainNet), // fixme
		Satoshi: int(txOut.Satoshi),

		GenesisHex:    hex.EncodeToString(txOut.Genesis),
		ScriptTypeHex: hex.EncodeToString(txOut.ScriptType),
		ScriptHex:     hex.EncodeToString(txOut.Script),
		Height:        int(txOut.Height),
	}
	return
}
