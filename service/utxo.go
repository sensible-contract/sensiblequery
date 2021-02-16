package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/lib/blkparser"
	"satoblock/lib/utils"
	"satoblock/model"
)

//////////////// address
func GetUtxoByAddress(addressHex string) (txOutsRsp []*model.TxOutResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM utxo_address
WHERE address = unhex('%s')
ORDER BY height DESC
LIMIT 32
`, addressHex)

	txOutsRet, err := clickhouse.ScanAll(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txs by address failed: %v", err)
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
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
SELECT txid, vout, address, genesis, satoshi, script_type, script, height FROM utxo_genesis
WHERE genesis = unhex('%s')
ORDER BY height DESC
LIMIT 32`,
		genesisHex)

	txOutsRet, err := clickhouse.ScanAll(psql, txOutResultSRF)
	if err != nil {
		log.Printf("query txs by genesis failed: %v", err)
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
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
