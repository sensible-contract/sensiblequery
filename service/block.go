package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/lib/blkparser"
	"satoblock/model"
)

const (
	SQL_FIELEDS_BLOCK = "height, blkid, previd, next_blk.blkid, merkle, ntx, invalue, outvalue, coinbase_out, blocktime, bits, blocksize"
)

func blockResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.BlockDO
	err := rows.Scan(&ret.Height, &ret.BlockId, &ret.PrevBlockId, &ret.NextBlockId, &ret.MerkleRoot, &ret.TxCount, &ret.InSatoshi, &ret.OutSatoshi, &ret.CoinbaseOut, &ret.BlockTime, &ret.Bits, &ret.BlockSize)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

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
		log.Printf("query blk failed: %v", err)
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

func GetBestBlock() (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blk_height ORDER BY height DESC LIMIT 1", SQL_FIELEDS_BLOCK)
	return GetBlockBySql(psql)
}

func GetBlockBySql(psql string) (blk *model.BlockInfoResp, err error) {
	blkRet, err := clickhouse.ScanOne(psql, blockResultSRF)
	if err != nil {
		log.Printf("query blk failed: %v", err)
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
