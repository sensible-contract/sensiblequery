package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/logger"
	"sensiblequery/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// "height, codehash, genesis, code_type, nft_idx, in_data_value, out_data_value, invalue, outvalue, blkid"
func nftInfoResultSRF(rows *sql.Rows) (interface{}, error) {
	var ret model.NFTInfoDO
	err := rows.Scan(&ret.CodeHash, &ret.Genesis, &ret.Count, &ret.InTimes, &ret.OutTimes, &ret.InSatoshi, &ret.OutSatoshi)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func getNFTMetaInfo(nftsRsp []*model.NFTInfoResp) {
	pipe := rdb.Pipeline()
	nftinfoCmds := make([]*redis.StringStringMapCmd, 0)
	for _, nft := range nftsRsp {
		// nftinfo of each token
		key, _ := hex.DecodeString(nft.CodeHashHex + nft.GenesisHex)
		nftinfoCmds = append(nftinfoCmds, pipe.HGetAll(ctx, "ni"+string(key)))
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}
	for idx, nft := range nftsRsp {
		nftinfo, err := nftinfoCmds[idx].Result()
		if err == redis.Nil {
			nftinfo = map[string]string{
				"metatxid":   "",
				"metavout":   "0",
				"supply":     "0",
				"sensibleid": "",
			}
			continue
		} else if err != nil {
			logger.Log.Info("getNFTDecimal redis failed", zap.Error(err))
		}
		supply, _ := strconv.Atoi(nftinfo["supply"])
		metavout, _ := strconv.Atoi(nftinfo["metavout"])
		nft.Supply = supply
		nft.MetaTxIdHex = hex.EncodeToString([]byte(nftinfo["metatxid"]))
		nft.MetaOutputIndex = metavout
		nft.SensibleIdHex = hex.EncodeToString([]byte(nftinfo["sensibleid"]))
	}
}

func ListNFTInfoByGenesis(codeHashHex, genesisHex string) (nftsRsp []*model.NFTInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT codehash, genesis, count(1), sum(in_times), sum(out_times), sum(in_satoshi), sum(out_satoshi) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times , sum(out_data_value) AS out_times,
            sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
     WHERE code_type = 3 AND codehash = unhex('%s') AND genesis = unhex('%s')
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`, codeHashHex, genesisHex)

	nftsRsp, err = GetNFTInfoBySQL(psql)
	if err != nil {
		return
	}
	getNFTMetaInfo(nftsRsp)
	return
}

func GetNFTSummary(codeHashHex string) (nftsRsp []*model.NFTInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT codehash, genesis, count(1), sum(in_times), sum(out_times), sum(in_satoshi), sum(out_satoshi) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times , sum(out_data_value) AS out_times,
            sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
     WHERE code_type = 3 AND codehash = unhex('%s')
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`, codeHashHex)

	nftsRsp, err = GetNFTInfoBySQL(psql)
	if err != nil {
		return
	}
	getNFTMetaInfo(nftsRsp)
	return
}

func GetNFTInfo() (nftsRsp []*model.NFTInfoResp, err error) {
	psql := `
SELECT codehash, genesis, count(1), sum(in_times), sum(out_times), sum(in_satoshi), sum(out_satoshi) FROM (
     SELECT codehash, genesis, nft_idx,
            sum(in_data_value) AS in_times , sum(out_data_value) AS out_times,
            sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
     WHERE code_type = 3
     GROUP BY codehash, genesis, nft_idx
)
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`
	nftsRsp, err = GetNFTInfoBySQL(psql)
	if err != nil {
		return
	}
	getNFTMetaInfo(nftsRsp)
	return
}

func GetNFTInfoBySQL(psql string) (blksRsp []*model.NFTInfoResp, err error) {
	blksRet, err := clickhouse.ScanAll(psql, nftInfoResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
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
