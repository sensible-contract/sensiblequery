package service

import (
	"database/sql"
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

func GetContractSwapDataInBlocksByHeightRange(cursor, size, blkStartHeight, blkEndHeight int, codeHashHex, genesisHex string, codeType uint32) (blksRsp []*model.ContractSwapDataResp, err error) {
	if blkEndHeight == 0 {
		blkEndHeight = 4294967295 + 1 // enable mempool
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
LIMIT %d, %d`,
		SQL_FIELEDS_SWAP_DATA,
		blkStartHeight, blkEndHeight,
		blkStartHeight, blkEndHeight, codeType, codeHashHex, genesisHex, cursor, size)

	blksRet, err := clickhouse.ScanAll(psql, contractSwapDataResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blksRet == nil {
		blksRsp = make([]*model.ContractSwapDataResp, 0)
		return
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

func contractSwapAggregateResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.ContractSwapAggregateDo
	err := rows.Scan(&ret.Height, &ret.BlockTime, &ret.OpenPrice, &ret.ClosePrice, &ret.MinPrice, &ret.MaxPrice, &ret.Token1Volume, &ret.Token2Volume)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetContractSwapAggregateInBlocksByHeightRange(interval, blkStartHeight, blkEndHeight int, codeHashHex, genesisHex string, codeType uint32) (blksRsp []*model.ContractSwapAggregateResp, err error) {
	if blkEndHeight == 0 {
		blkEndHeight = 4294967295 // disable mempool
	}
	if interval < 1 {
		interval = 1
	}
	if interval > 10000 {
		interval = 10000
	}

	psql := fmt.Sprintf(`
SELECT height, blocktime, open_price, close_price, min_price, max_price, volume1, volume2 FROM (
    SELECT ts * %d as height,
           anyLast(price) as open_price,
           any(price) as close_price,
           min(price) as min_price,
           max(price) as max_price,
           sum(volume1) as volume1,
           sum(volume2) as volume2
    FROM (
      SELECT intDiv(height, %d) as ts,
             out_value1 / out_value2 as price,
             abs(in_value1 - out_value1) as volume1,
             abs(in_value2 - out_value2) as volume2
      FROM blktx_contract_height
      WHERE height >= %d AND height < %d AND
        code_type = %d AND
        operation < 2 AND
         codehash = unhex('%s') AND
          genesis = unhex('%s')
      ORDER BY height DESC, txidx DESC
    )
    GROUP BY ts
) AS swap
LEFT JOIN (
    SELECT height, blocktime FROM blk_height
    WHERE height >= %d AND height < %d
) AS block
USING height
ORDER BY height DESC
`, interval, interval,
		blkStartHeight, blkEndHeight, codeType, codeHashHex, genesisHex,
		blkStartHeight, blkEndHeight)

	blksRet, err := clickhouse.ScanAll(psql, contractSwapAggregateResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blksRet == nil {
		blksRsp = make([]*model.ContractSwapAggregateResp, 0)
		return
	}
	blocks := blksRet.([]*model.ContractSwapAggregateDo)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.ContractSwapAggregateResp{
			Height:       int(block.Height),
			BlockTime:    int(block.BlockTime),
			OpenPrice:    block.OpenPrice,
			ClosePrice:   block.ClosePrice,
			MinPrice:     block.MinPrice,
			MaxPrice:     block.MaxPrice,
			Token1Volume: int(block.Token1Volume),
			Token2Volume: int(block.Token2Volume),
		})
	}
	return
}

////////////////////////////////////////////////////////////////

func contractSwapAggregateAmountResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.ContractSwapAggregateAmountDo
	err := rows.Scan(&ret.Height, &ret.BlockTime, &ret.OpenAmount, &ret.CloseAmount, &ret.MinAmount, &ret.MaxAmount)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetContractSwapAggregateAmountInBlocksByHeightRange(interval, blkStartHeight, blkEndHeight int, codeHashHex, genesisHex string, codeType uint32) (blksRsp []*model.ContractSwapAggregateAmountResp, err error) {
	if blkEndHeight == 0 {
		blkEndHeight = 4294967295 // disable mempool
	}
	if interval < 1 {
		interval = 1
	}
	if interval > 10000 {
		interval = 10000
	}

	psql := fmt.Sprintf(`
SELECT height, blocktime, open_amount, close_amount, min_amount, max_amount FROM (
    SELECT ts * %d as height,
           anyLast(amount) as open_amount,
           any(amount) as close_amount,
           min(amount) as min_amount,
           max(amount) as max_amount
    FROM (
      SELECT intDiv(height, %d) as ts,
             out_value1 as amount,
      FROM blktx_contract_height
      WHERE height >= %d AND height < %d AND
        code_type = %d AND
         codehash = unhex('%s') AND
          genesis = unhex('%s')
      ORDER BY height DESC, txidx DESC
    )
    GROUP BY ts
) AS swap
LEFT JOIN (
    SELECT height, blocktime FROM blk_height
    WHERE height >= %d AND height < %d
) AS block
USING height
ORDER BY height DESC
`, interval, interval,
		blkStartHeight, blkEndHeight, codeType, codeHashHex, genesisHex,
		blkStartHeight, blkEndHeight)

	blksRet, err := clickhouse.ScanAll(psql, contractSwapAggregateAmountResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
		return nil, err
	}
	if blksRet == nil {
		blksRsp = make([]*model.ContractSwapAggregateAmountResp, 0)
		return
	}
	blocks := blksRet.([]*model.ContractSwapAggregateAmountDo)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.ContractSwapAggregateAmountResp{
			Height:      int(block.Height),
			BlockTime:   int(block.BlockTime),
			OpenAmount:  int(block.OpenAmount),
			CloseAmount: int(block.CloseAmount),
			MinAmount:   int(block.MinAmount),
			MaxAmount:   int(block.MaxAmount),
		})
	}
	return
}
