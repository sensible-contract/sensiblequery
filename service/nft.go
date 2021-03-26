package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"satoblock/dao/clickhouse"
	"satoblock/model"
)

const (
	SQL_FIELEDS_NFT_VOLUME_INFO = "height, codehash, genesis, code_type, nft_idx, in_data_value, out_data_value, invalue, outvalue, blkid"
)

func nftInfoResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.NFTInfoDO
	err := rows.Scan(&ret.CodeHash, &ret.Genesis, &ret.Count, &ret.InTimes, &ret.OutTimes, &ret.InSatoshi, &ret.OutSatoshi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetNFTInfo() (blksRsp []*model.NFTInfoResp, err error) {
	psql := `
SELECT codehash, genesis, count(nft_idx), sum(in_times), sum(out_times), sum(in_satoshi), sum(out_satoshi) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times , sum(out_data_value) AS out_times,
            sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
     WHERE code_type = 0
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash, genesis
`

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
