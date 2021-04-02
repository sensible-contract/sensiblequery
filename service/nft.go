package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/model"
)

// "height, codehash, genesis, code_type, nft_idx, in_data_value, out_data_value, invalue, outvalue, blkid"

func tokenCodeHashResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.TokenCodeHashDO
	err := rows.Scan(&ret.CodeHash, &ret.Count, &ret.InTimes, &ret.OutTimes)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTokenCodeHash(codeType int) (blksRsp []*model.TokenCodeHashResp, err error) {
	psql := fmt.Sprintf(`
SELECT codehash, count(1), sum(in_times), sum(out_times) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times, sum(out_data_value) AS out_times FROM blk_codehash_height
     WHERE code_type = %d
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash
ORDER BY count(1) DESC
`, codeType)
	blksRet, err := clickhouse.ScanAll(psql, tokenCodeHashResultSRF)
	if err != nil {
		log.Printf("query nft codehash failed: %v", err)
		return nil, err
	}
	if blksRet == nil {
		return nil, errors.New("not exist")
	}
	blocks := blksRet.([]*model.TokenCodeHashDO)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.TokenCodeHashResp{
			CodeHashHex: hex.EncodeToString(block.CodeHash),
			Count:       int(block.Count),
			InTimes:     int(block.InTimes),
			OutTimes:    int(block.OutTimes),
		})
	}
	return
}

func nftInfoResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.NFTInfoDO
	err := rows.Scan(&ret.CodeHash, &ret.Genesis, &ret.Count, &ret.InTimes, &ret.OutTimes, &ret.InSatoshi, &ret.OutSatoshi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetNFTSummary(codeHashHex string) (blksRsp []*model.NFTInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT codehash, genesis, count(1), sum(in_times), sum(out_times), sum(in_satoshi), sum(out_satoshi) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times , sum(out_data_value) AS out_times,
            sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
     WHERE code_type = 0 AND codehash = unhex('%s')
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`, codeHashHex)

	return GetNFTInfoBySQL(psql)
}

func GetNFTInfo() (blksRsp []*model.NFTInfoResp, err error) {
	psql := `
SELECT codehash, genesis, count(1), sum(in_times), sum(out_times), sum(in_satoshi), sum(out_satoshi) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times , sum(out_data_value) AS out_times,
            sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
     WHERE code_type = 0
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`
	return GetNFTInfoBySQL(psql)
}

func GetNFTInfoBySQL(psql string) (blksRsp []*model.NFTInfoResp, err error) {
	blksRet, err := clickhouse.ScanAll(psql, nftInfoResultSRF)
	if err != nil {
		log.Printf("query blk failed: %v", err)
		return nil, err
	}
	if blksRet == nil {
		return nil, errors.New("not exist")
	}
	blocks := blksRet.([]*model.NFTInfoDO)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.NFTInfoResp{
			CodeHashHex: hex.EncodeToString(block.CodeHash),
			GenesisHex:  hex.EncodeToString(block.Genesis),
			Count:       int(block.Count),
			InTimes:     int(block.InTimes),
			OutTimes:    int(block.OutTimes),
			InSatoshi:   int(block.InSatoshi),
			OutSatoshi:  int(block.OutSatoshi),
		})
	}
	return

}
