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
	SQL_FIELEDS_FT_VOLUME_INFO = "height, codehash, genesis, code_type, nft_idx, in_data_value, out_data_value, invalue, outvalue, blkid"
)

func ftInfoResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.FTInfoDO
	err := rows.Scan(&ret.CodeHash, &ret.Genesis, &ret.Count, &ret.InVolume, &ret.OutVolume, &ret.InSatoshi, &ret.OutSatoshi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetFTInfo() (blksRsp []*model.FTInfoResp, err error) {
	psql := `
SELECT codehash, genesis, count(1),
       sum(in_data_value) AS in_volume , sum(out_data_value) AS out_volume,
       sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
WHERE code_type = 1
GROUP BY codehash, genesis
`

	blksRet, err := clickhouse.ScanAll(psql, ftInfoResultSRF)
	if err != nil {
		log.Printf("query blk failed: %v", err)
		return nil, err
	}
	if blksRet == nil {
		return nil, errors.New("not exist")
	}
	blocks := blksRet.([]*model.FTInfoDO)
	for _, block := range blocks {
		blksRsp = append(blksRsp, &model.FTInfoResp{
			CodeHashHex: hex.EncodeToString(block.CodeHash),
			GenesisHex:  hex.EncodeToString(block.Genesis),
			Count:       int(block.Count),
			InVolume:    int(block.InVolume),
			OutVolume:   int(block.OutVolume),
			InSatoshi:   int(block.InSatoshi),
			OutSatoshi:  int(block.OutSatoshi),
		})
	}
	return

}
