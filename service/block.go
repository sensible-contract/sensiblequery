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
	SQL_FIELEDS_BLOCK = "height, blkid, previd, merkle, ntx, blocktime, bits, blocksize"
)

func blockResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.BlockDO
	err := rows.Scan(&ret.Height, &ret.BlockId, &ret.PrevBlockId, &ret.MerkleRoot, &ret.TxCount, &ret.BlockTime, &ret.Bits, &ret.BlockSize)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetBlocksByHeightRange(blkStartHeight, blkEndHeight int) (blksRsp []*model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blk_height WHERE height >= %d AND height < %d",
		SQL_FIELEDS_BLOCK, blkStartHeight, blkEndHeight)

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
			MerkleRootHex:  blkparser.HashString(block.MerkleRoot),
			TxCount:        int(block.TxCount),
			BlockTime:      int(block.BlockTime),
			Bits:           int(block.Bits),
			BlockSize:      int(block.BlockSize),
		})
	}
	return

}

func GetBlockByHeight(blkHeight int) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blk_height WHERE height = %d", SQL_FIELEDS_BLOCK, blkHeight)
	return GetBlockBySql(psql)
}

func GetBlockById(blkidHex string) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT %s FROM blk WHERE blkid = unhex('%s')", SQL_FIELEDS_BLOCK, blkidHex)
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
		MerkleRootHex:  blkparser.HashString(block.MerkleRoot),
		TxCount:        int(block.TxCount),
		BlockTime:      int(block.BlockTime),
		Bits:           int(block.Bits),
		BlockSize:      int(block.BlockSize),
	}
	return
}
