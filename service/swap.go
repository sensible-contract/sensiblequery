package service

import (
	"database/sql"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/logger"
	"sensiblequery/model"

	"go.uber.org/zap"
)

const (
	SQL_FIELEDS_SWAP_DATA = "height, blk.blocktime, code_type, operation, in_value1, in_value2, in_value3, out_value1, out_value2, out_value3, txidx"
)

func contractSwapDataResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.ContractSwapDataDo
	err := rows.Scan(&ret.Height, &ret.BlockTime, &ret.CodeType, &ret.Operation, &ret.InToken1Amount, &ret.InToken2Amount, &ret.InLpAmount, &ret.OutToken1Amount, &ret.OutToken2Amount, &ret.OutLpAmount, &ret.Idx)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetContractSwapDataInBlocksByHeightRange(size, blkStartHeight, blkEndHeight int, codeHashHex, genesisHex string, codeType uint32) (blksRsp []*model.ContractSwapDataResp, err error) {
	if blkEndHeight == 0 {
		blkEndHeight = 4294967296
	}
	psql := fmt.Sprintf(`
SELECT %s FROM blktx_contract_height
LEFT JOIN (
    SELECT height, blocktime FROM blk_height
    WHERE height >= %d AND height < %d
) AS blk
USING height
WHERE height >= %d AND height < %d AND
    code_type = %d AND
     codehash = unhex('%s') AND
      genesis = unhex('%s')
ORDER BY height DESC, txidx DESC
LIMIT %d`,
		SQL_FIELEDS_SWAP_DATA,
		blkStartHeight, blkEndHeight,
		blkStartHeight, blkEndHeight, codeType, codeHashHex, genesisHex, size)

	blksRet, err := clickhouse.ScanAll(psql, contractSwapDataResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blksRet == nil {
		return nil, errors.New("not exist")
	}
	blocks := blksRet.([]*model.ContractSwapDataDo)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.ContractSwapDataResp{
			Height:          int(block.Height),
			BlockTime:       int(block.BlockTime),
			CodeType:        int(block.CodeType),
			Operation:       int(block.Operation),
			InToken1Amount:  int(block.InToken1Amount),
			InToken2Amount:  int(block.InToken2Amount),
			InLpAmount:      int(block.InLpAmount),
			OutToken1Amount: int(block.OutToken1Amount),
			OutToken2Amount: int(block.OutToken2Amount),
			OutLpAmount:     int(block.OutLpAmount),
			Idx:             int(block.Idx),
		})
	}
	return

}
