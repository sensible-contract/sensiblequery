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

//////////////// history
func txOutHistoryResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TxOutHistoryDO
	err := rows.Scan(&ret.TxId, &ret.Vout, &ret.Address, &ret.Genesis, &ret.Satoshi, &ret.ScriptType, &ret.ScriptPk, &ret.Height, &ret.Idx, &ret.IOType)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//////////////// address
func GetHistoryByAddress(addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	psql := fmt.Sprintf(`
SELECT txid, idx, address, genesis, satoshi, script_type, script_pk, height, txidx, io_type FROM
(
SELECT utxid AS txid, vout AS idx, address, genesis, satoshi, script_type, script_pk, height, utxidx AS txidx, 1 AS io_type FROM txout
WHERE (utxid, vout, height) in (
    SELECT utxid, vout, height FROM txout_address_height
    WHERE address = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)

UNION ALL

SELECT txid, idx, address, genesis, satoshi, script_type, script_pk, height, txidx, 0 AS io_type FROM txin
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_address_height
    WHERE address = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)
)
ORDER BY height DESC, txidx DESC
LIMIT 128
`, addressHex, addressHex)
	return GetHistoryBySql(psql)
}

//////////////// genesis
func GetHistoryByGenesis(cursor, size int, codeHashHex, genesisHex, addressHex string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	logger.Log.Info("query tx history by codehash/genesis for", zap.String("address", addressHex))
	psql := fmt.Sprintf(`
SELECT txid, idx, address, genesis, satoshi, script_type, script_pk, height, txidx, io_type FROM
(
SELECT utxid AS txid, vout AS idx, address, genesis, satoshi, script_type, script_pk, height, utxidx AS txidx, 1 AS io_type FROM txout
WHERE (utxid, vout, height) in (
    SELECT utxid, vout, height FROM txout_genesis_height
    WHERE codehash = unhex('%s') AND
          genesis = unhex('%s') AND
          address = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)

UNION ALL

SELECT txid, idx, address, genesis, satoshi, script_type, script_pk, height, txidx, 0 AS io_type FROM txin
WHERE (txid, idx, height) in (
    SELECT txid, idx, height FROM txin_genesis_height
    WHERE codehash = unhex('%s') AND
          genesis = unhex('%s') AND
          address = unhex('%s')
    ORDER BY height DESC
    LIMIT 64
)
)
ORDER BY height DESC, txidx DESC
LIMIT 128
`, codeHashHex, genesisHex, addressHex,
		codeHashHex, genesisHex, addressHex)
	return GetHistoryBySql(psql)
}

func GetHistoryBySql(psql string) (txOutsRsp []*model.TxOutHistoryResp, err error) {
	txOutsRet, err := clickhouse.ScanAll(psql, txOutHistoryResultSRF)
	if err != nil {
		logger.Log.Info("query tx history by genesis failed", zap.Error(err))
		return nil, err
	}
	if txOutsRet == nil {
		return nil, errors.New("not exist")
	}
	txOuts := txOutsRet.([]*model.TxOutHistoryDO)
	for _, txout := range txOuts {

		txo := scriptDecoder.ExtractPkScriptForTxo(txout.ScriptPk, txout.ScriptType)

		txOutsRsp = append(txOutsRsp, &model.TxOutHistoryResp{
			TxIdHex: blkparser.HashString(txout.TxId),
			Vout:    int(txout.Vout),
			Address: utils.EncodeAddress(txout.Address, utils.PubKeyHashAddrID),
			Satoshi: int(txout.Satoshi),

			IsNFT:           (txo.CodeType == scriptDecoder.CodeType_NFT),
			CodeType:        int(txo.CodeType),
			TokenIndex:      strconv.FormatUint(txo.TokenIndex, 10),
			MetaTxIdHex:     hex.EncodeToString(txo.MetaTxId),
			MetaOutputIndex: int(txo.MetaOutputIndex),
			TokenId:         hex.EncodeToString(txo.GenesisId),
			TokenName:       txo.Name,
			TokenSymbol:     txo.Symbol,
			TokenAmount:     strconv.FormatUint(txo.Amount, 10),
			TokenDecimal:    int(txo.Decimal),
			CodeHashHex:     hex.EncodeToString(txo.CodeHash),
			GenesisHex:      hex.EncodeToString(txo.GenesisId),
			ScriptTypeHex:   hex.EncodeToString(txout.ScriptType),

			Height: int(txout.Height),
			Idx:    int(txout.Idx),
			IOType: int(txout.IOType),
		})
	}
	return
}
