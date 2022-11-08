package service

import (
	"encoding/hex"
	"fmt"
	"sensiblequery/dao/rdb"
	"sensiblequery/logger"
	"sensiblequery/model"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

////////////////
func GetTxsHistoryInfoByAddress(addressPkh []byte) (addrRsp *model.AddressHistoryInfoResp, err error) {
	historyNum, err := rdb.RdbAddressClient.ZCard(ctx, "{ah"+string(addressPkh)+"}").Result()
	if err != nil {
		logger.Log.Info("get historyNum from redis failed", zap.Error(err))
		return
	}
	logger.Log.Info("historyNum", zap.Int64("n", historyNum))

	addrRsp = &model.AddressHistoryInfoResp{
		Total: int(historyNum),
	}
	return addrRsp, nil
}

func GetTxsHistoryByAddressAndTypeByHeightRangeFromPika(cursor, size int, addressPkh []byte) (txsRsp []*model.TxInfoResp, err error) {
	addrTxWithHeightHistory, err := rdb.RdbAddressClient.ZRevRange(ctx, "{ah"+string(addressPkh)+"}", int64(cursor), int64(cursor+size)-1).Result()
	if err == redis.Nil {
		addrTxWithHeightHistory = nil
	} else if err != nil {
		logger.Log.Info("GetTxsHistoryByAddressAndTypeByHeightRangeFromPika failed", zap.Error(err))
		return
	}

	for _, historyPosition := range addrTxWithHeightHistory {
		sep := strings.Index(historyPosition, ":")
		height, _ := strconv.Atoi(historyPosition[:sep])
		txidx, _ := strconv.Atoi(historyPosition[sep+1:])
		txsRsp = append(txsRsp, &model.TxInfoResp{
			Height: height,
			Idx:    txidx,
		})
	}

	return txsRsp, nil
}

//////////////// address
func GetTxsHistoryByAddressAndTypeByHeightRange(cursor, size int, addressPkh []byte, historyType model.HistoryType) (txsRsp []*model.TxInfoResp, err error) {
	logger.Log.Info("query txinfo history for",
		zap.Int("cursor", cursor),
		zap.Int("size", size),
		zap.String("address", hex.EncodeToString(addressPkh)))

	txsRsp, err = GetTxsHistoryByAddressAndTypeByHeightRangeFromPika(cursor, size, addressPkh)
	if err != nil || len(txsRsp) == 0 {
		return
	}

	blkStartHeight := 4294967295
	blkEndHeight := 0
	strHeightTxidList := make([]string, len(txsRsp))
	for idx, tx := range txsRsp {
		if blkStartHeight > tx.Height {
			blkStartHeight = tx.Height
		}
		if blkEndHeight < tx.Height {
			blkEndHeight = tx.Height
		}
		strHeightTxidList[idx] = fmt.Sprintf("(%d,%d)", tx.Height, tx.Idx)
	}

	psql := fmt.Sprintf(`
SELECT %s FROM blktx_height
LEFT JOIN  (
    SELECT height, blkid, blocktime FROM blk_height
    WHERE height >= %d AND height <= %d
) AS blk
USING height
WHERE (height, txidx) in (%s)
ORDER BY height DESC, txidx DESC
`,
		SQL_FIELEDS_TX_TIMESTAMP,
		blkStartHeight, blkEndHeight,

		strings.Join(strHeightTxidList, ","))

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
)
ORDER BY height DESC, txidx DESC
LIMIT %d, %d
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
