package service

import (
	"fmt"
	"sensiblequery/logger"
	"sensiblequery/model"

	"go.uber.org/zap"
)

//////////////// address
func GetTxsHistoryByAddressAndTypeByHeightRange(cursor, size, blkStartHeight, blkEndHeight int, addressHex string, historyType model.HistoryType) (txsRsp []*model.TxInfoResp, err error) {
	logger.Log.Info("query txinfo history for",
		zap.Int("cursor", cursor),
		zap.Int("size", size),
		zap.Int("blkStart", blkStartHeight),
		zap.Int("blkEnd", blkEndHeight),
		zap.String("address", addressHex))
	maxOffset := cursor + size

	if blkEndHeight == 0 {
		blkEndHeight = 4294967295 + 1 // enable mempool
	}

	codehashMatch := ""
	if historyType == model.HISTORY_CONTRACT_ONLY {
		codehashMatch = "AND codehash != '' AND codehash != unhex('00')"
	}

	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
LEFT JOIN  (
    SELECT height, blkid, blocktime FROM blk_height
    WHERE height >= %d AND height < %d
) AS blk
USING height
WHERE (height, txidx, substring(txid, 1, 12)) in (
    SELECT height, txidx, txid FROM
    (
        SELECT utxid AS txid, height, utxidx AS txidx FROM txout_address_height
        WHERE height >= %d AND height < %d AND address = unhex('%s') %s
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC
        LIMIT %d

      UNION ALL

        SELECT substring(utxid, 1, 12) AS txid, height, utxidx AS txidx FROM txout
        WHERE height >= 4294967295 AND height < %d AND address = unhex('%s') %s
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC
        LIMIT %d

      UNION ALL

        SELECT txid, height, txidx FROM txin_address_height
        WHERE height >= %d AND height < %d AND address = unhex('%s') %s
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC
        LIMIT %d

      UNION ALL

        SELECT substring(txid, 1, 12), height, txidx FROM txin
        WHERE height >= 4294967295 AND height < %d AND address = unhex('%s') %s
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC
        LIMIT %d
    ) AS txlist
    ORDER BY height DESC, txidx DESC
    LIMIT %d, %d
)
ORDER BY height DESC, txidx DESC
`,
		SQL_FIELEDS_TX_TIMESTAMP,
		blkStartHeight, blkEndHeight,

		blkStartHeight, blkEndHeight,
		addressHex, codehashMatch, maxOffset,
		blkEndHeight,
		addressHex, codehashMatch, maxOffset,

		blkStartHeight, blkEndHeight,
		addressHex, codehashMatch, maxOffset,
		blkEndHeight,
		addressHex, codehashMatch, maxOffset,

		cursor, size)

	return GetBlockTxsBySql(psql, true)
}

//////////////// genesis
func GetTxsHistoryByGenesisByHeightRange(cursor, size, blkStartHeight, blkEndHeight int, codehashHex, genesisHex, addressHex string) (txsRsp []*model.TxInfoResp, err error) {
	logger.Log.Info("query txinfo history by codehash/genesis for",
		zap.Int("cursor", cursor),
		zap.Int("size", size),
		zap.Int("blkStart", blkStartHeight),
		zap.Int("blkEnd", blkEndHeight),
		zap.String("address", addressHex))

	if blkEndHeight == 0 {
		blkEndHeight = 4294967295 + 1 // enable mempool
	}

	addressMatch := ""
	if addressHex == "0000000000000000000000000000000000000000" {
		addressMatch = "OR address = ''"
	}

	codehashMatch := ""
	genesisMatch := ""
	if codehashHex != "0000000000000000000000000000000000000000" {
		codehashMatch = fmt.Sprintf("codehash = unhex('%s') AND", codehashHex)
		genesisMatch = fmt.Sprintf("genesis = unhex('%s') AND", genesisHex)
	}
	maxOffset := cursor + size

	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
LEFT JOIN  (
    SELECT height, blkid, blocktime FROM blk_height
    WHERE height >= %d AND height < %d
) AS blk
USING height
WHERE (substring(txid, 1, 12), height, txidx) in (
    SELECT txid, height, txidx FROM
    (
        SELECT utxid AS txid, height, utxidx AS txidx FROM txout_genesis_height
        WHERE height >= %d AND height < %d AND %s %s (address = unhex('%s') %s)
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC, codehash DESC, genesis DESC
        LIMIT %d

      UNION ALL

        SELECT substring(utxid, 1, 12) AS txid, height, utxidx AS txidx FROM txout
        WHERE height >= 4294967295 AND height < %d AND %s %s (address = unhex('%s') %s)
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC, codehash DESC, genesis DESC
        LIMIT %d

      UNION ALL

        SELECT txid, height, txidx FROM txin_genesis_height
        WHERE height >= %d AND height < %d AND %s %s (address = unhex('%s') %s)
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC, codehash DESC, genesis DESC
        LIMIT %d

      UNION ALL

        SELECT substring(txid, 1, 12), height, txidx FROM txin
        WHERE height >= 4294967295 AND height < %d AND %s %s (address = unhex('%s') %s)
        GROUP BY txid, height, txidx
        ORDER BY height DESC, txidx DESC, codehash DESC, genesis DESC
        LIMIT %d
    ) AS txlist
    ORDER BY height DESC, txidx DESC
    LIMIT %d, %d
)
ORDER BY height DESC, txidx DESC
`,
		SQL_FIELEDS_TX_TIMESTAMP,
		blkStartHeight, blkEndHeight,

		blkStartHeight, blkEndHeight,
		codehashMatch, genesisMatch, addressHex, addressMatch, maxOffset,
		blkEndHeight,
		codehashMatch, genesisMatch, addressHex, addressMatch, maxOffset,

		blkStartHeight, blkEndHeight,
		codehashMatch, genesisMatch, addressHex, addressMatch, maxOffset,
		blkEndHeight,
		codehashMatch, genesisMatch, addressHex, addressMatch, maxOffset,

		cursor, size)

	return GetBlockTxsBySql(psql, true)
}
