package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/dao/rdb"
	"sensiblequery/lib/blkparser"
	"sensiblequery/logger"
	"sensiblequery/model"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	SQL_FIELEDS_BEST_BLOCK   = "height, blkid, previd, '0000000000000000000000000000000000000000000000000000000000000000', merkle, ntx, invalue, outvalue, coinbase_out, blocktime, bits, blocksize"
	SQL_FIELEDS_BLOCK        = "height, blkid, previd, next_blk.blkid, merkle, ntx, invalue, outvalue, coinbase_out, blocktime, bits, blocksize"
	SQL_FIELEDS_BLOCK_VOLUME = "height, codehash, genesis, code_type, nft_idx, in_data_value, out_data_value, invalue, outvalue, blkid"
)

func blockTokenVolumeResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.BlockTokenVolumeDO
	err := rows.Scan(&ret.Height, &ret.CodeHash, &ret.Genesis, &ret.CodeType, &ret.NFTIdx, &ret.InDataValue, &ret.OutDataValue, &ret.InSatoshi, &ret.OutSatoshi, &ret.BlockId)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func blockResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.BlockDO
	err := rows.Scan(&ret.Height, &ret.BlockId, &ret.PrevBlockId, &ret.NextBlockId, &ret.MerkleRoot, &ret.TxCount, &ret.InSatoshi, &ret.OutSatoshi, &ret.CoinbaseOut, &ret.BlockTime, &ret.Bits, &ret.BlockSize)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTokenVolumesInBlocksByHeightRange(blkStartHeight, blkEndHeight int, codeHashHex, genesisHex string, codeType uint32, nftIdx int) (blksRsp []*model.BlockTokenVolumeResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blk_codehash_height
WHERE height >= %d AND height < %d AND
   codehash = unhex('%s') AND
    genesis = unhex('%s') AND
    code_type = %d AND
    nft_idx = %d
ORDER BY height ASC
LIMIT %d`,
		SQL_FIELEDS_BLOCK_VOLUME, blkStartHeight, blkEndHeight, codeHashHex, genesisHex, codeType, nftIdx, blkEndHeight-blkStartHeight)

	blksRet, err := clickhouse.ScanAll(psql, blockTokenVolumeResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blksRet == nil {
		return nil, errors.New("not exist")
	}
	blocks := blksRet.([]*model.BlockTokenVolumeDO)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.BlockTokenVolumeResp{
			Height:       int(block.Height),
			CodeHashHex:  hex.EncodeToString(block.CodeHash),
			GenesisHex:   hex.EncodeToString(block.Genesis),
			CodeType:     int(block.CodeType),
			NFTIdx:       int(block.NFTIdx),
			InDataValue:  int(block.InDataValue),
			OutDataValue: int(block.OutDataValue),
			InSatoshi:    int(block.InSatoshi),
			OutSatoshi:   int(block.OutSatoshi),
			BlockIdHex:   blkparser.HashString(block.BlockId),
		})
	}
	return

}

////////////////////////////////////////////////////////////////
func GetBlocksByHeightRange(blkStartHeight, blkEndHeight int) (blksRsp []*model.BlockInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blk_height
LEFT JOIN (
    SELECT blkid, previd FROM blk_height
    WHERE height > %d AND height <= %d
    LIMIT %d
) AS next_blk
ON blk_height.blkid = next_blk.previd
WHERE height >= %d AND height < %d ORDER BY height ASC
LIMIT %d
`,
		SQL_FIELEDS_BLOCK, blkStartHeight, blkEndHeight, blkEndHeight-blkStartHeight,
		blkStartHeight, blkEndHeight, blkEndHeight-blkStartHeight)

	blksRet, err := clickhouse.ScanAll(psql, blockResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blksRet == nil {
		return nil, errors.New("not exist")
	}
	blocks := blksRet.([]*model.BlockDO)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.BlockInfoResp{
			Height:         int(block.Height),
			BlockIdHex:     blkparser.HashString(block.BlockId),
			PrevBlockIdHex: blkparser.HashString(block.PrevBlockId),
			NextBlockIdHex: blkparser.HashString(block.NextBlockId),
			MerkleRootHex:  blkparser.HashString(block.MerkleRoot),
			TxCount:        int(block.TxCount),
			InSatoshi:      int(block.InSatoshi),
			OutSatoshi:     int(block.OutSatoshi),
			CoinbaseOut:    int(block.CoinbaseOut),

			BlockTime: int(block.BlockTime),
			Bits:      int(block.Bits),
			BlockSize: int(block.BlockSize),
		})
	}
	return

}

func GetBlockByHeight(blkHeight int) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT %s FROM blk_height
LEFT JOIN (
    SELECT blkid, previd FROM blk_height
    WHERE height = %d+1
    LIMIT 1
) AS next_blk
ON blk_height.blkid = next_blk.previd
WHERE height = %d ORDER BY height ASC
LIMIT 1`, SQL_FIELEDS_BLOCK, blkHeight, blkHeight)
	return GetBlockBySql(psql)
}

func GetBlockById(blkidHex string) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf(`SELECT %s FROM blk
LEFT JOIN (
    SELECT blkid, previd FROM blk_height
    WHERE height IN (
       SELECT toUInt32(height+1) FROM blk
       WHERE blkid = unhex('%s')
       LIMIT 1
    )
) AS next_blk
ON blk.blkid = next_blk.previd
WHERE blkid = unhex('%s')
LIMIT 1`, SQL_FIELEDS_BLOCK, blkidHex, blkidHex)
	return GetBlockBySql(psql)
}

func GetBestBlockByHeight(blkHeight int) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blk_height WHERE height = %d LIMIT 1", SQL_FIELEDS_BEST_BLOCK, blkHeight)
	return GetBlockBySql(psql)
}

func GetBlockBySql(psql string) (blk *model.BlockInfoResp, err error) {
	blkRet, err := clickhouse.ScanOne(psql, blockResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blkRet == nil {
		return nil, errors.New("not exist")
	}
	block := blkRet.(*model.BlockDO)
	blk = &model.BlockInfoResp{
		Height:         int(block.Height),
		BlockIdHex:     blkparser.HashString(block.BlockId),
		PrevBlockIdHex: blkparser.HashString(block.PrevBlockId),
		NextBlockIdHex: blkparser.HashString(block.NextBlockId),
		MerkleRootHex:  blkparser.HashString(block.MerkleRoot),
		TxCount:        int(block.TxCount),
		InSatoshi:      int(block.InSatoshi),
		OutSatoshi:     int(block.OutSatoshi),
		CoinbaseOut:    int(block.CoinbaseOut),

		BlockTime: int(block.BlockTime),
		Bits:      int(block.Bits),
		BlockSize: int(block.BlockSize),
	}
	return
}

func mempoolResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret int
	err := rows.Scan(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetMempoolTxCount() (count int, err error) {
	psql := "SELECT count(1) FROM blktx_height WHERE height >= 4294967295"

	blkRet, err := clickhouse.ScanOne(psql, mempoolResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return 0, err
	}
	if blkRet == nil {
		return 0, errors.New("not exist")
	}
	mempoolTxCount := blkRet.(int)

	return mempoolTxCount, nil
}

func GetBlockMedianTimePast(height int) (mtp int, err error) {
	psql := fmt.Sprintf(`
SELECT toUInt32(quantileExact(blocktime)) FROM (
    SELECT blocktime FROM blk_height WHERE height > %d AND height <= %d
)
`, height-11, height)

	blkRet, err := clickhouse.ScanOne(psql, mempoolResultSRF)
	if err != nil {
		logger.Log.Info("query mtp failed", zap.Error(err))
		return 0, err
	}
	if blkRet == nil {
		return 0, errors.New("not exist")
	}
	mtp = blkRet.(int)

	return mtp, nil
}

////////////////
func GetBestBlockHeight() (height int, err error) {
	// get decimal from f info
	height, err = rdb.BizClient.HGet(ctx, "info", "blocks_total").Int()
	if err == redis.Nil {
		height = 0
		logger.Log.Info("GetBestBlockHeight, but info missing")
	} else if err != nil {
		logger.Log.Info("GetBestBlockHeight, but redis failed", zap.Error(err))
		return
	}

	return height, nil
}
