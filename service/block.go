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

func blockResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.BlockDO
	err := rows.Scan(&ret.Height, &ret.BlockId, &ret.PrevBlockId, &ret.TxCount)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetBlocksByHeightRange(blkStartHeight, blkEndHeight int) (blksRsp []*model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT height, blkid, previd, ntx FROM blk_height WHERE height >= %d AND height < %d",
		blkStartHeight, blkEndHeight)

	blksRet, err := clickhouse.ScanAll(psql, blockResultSRF)
	if err != nil {
		log.Printf("query blk failed: %v", err)
		return nil, err
	}

	blocks := blksRet.([]*model.BlockDO)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.BlockInfoResp{
			Height:         int(block.Height),
			BlockIdHex:     blkparser.HashString(block.BlockId),
			PrevBlockIdHex: blkparser.HashString(block.PrevBlockId),
			TxCount:        int(block.TxCount),
		})
	}
	return

}

func GetBlockByHeight(blkHeight int) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT height, blkid, previd, ntx FROM blk_height WHERE height = %d", blkHeight)
	return GetBlockBySql(psql)
}

func GetBlockById(blkidHex string) (blk *model.BlockInfoResp, err error) {
	psql := fmt.Sprintf("SELECT height, blkid, previd, ntx FROM blk WHERE blkid = unhex('%s')", blkidHex)
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
		TxCount:        int(block.TxCount),
	}
	return
}
