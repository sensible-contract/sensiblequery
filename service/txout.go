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
func txOutResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.Script, &ret.Height)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTxOutputsByTxId(blkHeight int, txidHex string) (txOutsRsp []*model.TxOutResp, err error) {
	var psql string
	if blkHeight < 0 {
		psql = fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout
WHERE txid = unhex('%s') AND
height IN (
    SELECT height FROM tx
    WHERE txid = unhex('%s') LIMIT 1
)`, txidHex, txidHex)
	} else {
		psql = fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout
WHERE txid = unhex('%s') AND
    height = %d`, txidHex, blkHeight)
	}

	txOutsRet, err := clickhouse.ScanAll(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txs by blkid failed: %v", err)
		return nil, err
	}
	txOuts := txOutsRet.([]*model.TxOutDO)
	for _, txout := range txOuts {
		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet), // fixme
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptHex:     hex.EncodeToString(txout.Script),
			Height:        int(txout.Height),
		})
	}
	return
}

func GetTxOutputByTxIdAndIdx(blkHeight int, txidHex string, index int) (txOutRsp *model.TxOutResp, err error) {
	var psql string
	if blkHeight < 0 {
		psql = fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout
WHERE txid = unhex('%s') AND
      vout = %d
LIMIT 1`, txidHex, index)
	} else {
		psql = fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout
WHERE txid = unhex('%s') AND
      vout = %d AND
    height = %d
LIMIT 1`, txidHex, index, blkHeight)
	}

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

//////////////// address
func GetUtxoByAddress(addressHex string) (txOutsRsp []*model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout_address
WHERE address = unhex('%s')
ORDER BY height DESC
LIMIT 100`,
		addressHex)

	txOutsRet, err := clickhouse.ScanAll(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txs by address failed: %v", err)
		return nil, err
	}
	txOuts := txOutsRet.([]*model.TxOutDO)
	for _, txout := range txOuts {
		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet), // fixme
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptHex:     hex.EncodeToString(txout.Script),
			Height:        int(txout.Height),
		})
	}
	return
}

//////////////// genesis
func GetUtxoByGenesis(genesisHex string) (txOutsRsp []*model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM txout_genesis
WHERE genesis = unhex('%s')
ORDER BY height DESC
LIMIT 100`,
		genesisHex)

	txOutsRet, err := clickhouse.ScanAll(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txs by genesis failed: %v", err)
		return nil, err
	}
	txOuts := txOutsRet.([]*model.TxOutDO)
	for _, txout := range txOuts {
		txOutsRsp = append(txOutsRsp, &model.TxOutResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrIDMainNet), // fixme
			Satoshi: int(txout.Satoshi),

			GenesisHex:    hex.EncodeToString(txout.Genesis),
			ScriptTypeHex: hex.EncodeToString(txout.ScriptType),
			ScriptHex:     hex.EncodeToString(txout.Script),
			Height:        int(txout.Height),
		})
	}
	return
}
