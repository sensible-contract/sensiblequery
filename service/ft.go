package service

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"sensiblequery/dao/clickhouse"
	"sensiblequery/dao/rdb"
	"sensiblequery/logger"
	"sensiblequery/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
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

func getFTDecimal(ftsRsp []*model.FTInfoResp) {
	pipe := rdb.BizClient.Pipeline()
	ftinfoCmds := make([]*redis.StringStringMapCmd, 0)
	for _, ft := range ftsRsp {
		// ftinfo of each token
		key, _ := hex.DecodeString(ft.CodeHashHex + ft.GenesisHex)
		ftinfoCmds = append(ftinfoCmds, pipe.HGetAll(ctx, "fi"+string(key)))
	}
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}
	for idx, ft := range ftsRsp {
		ftinfo, err := ftinfoCmds[idx].Result()
		if err == redis.Nil {
			ftinfo = map[string]string{
				"decimal":    "0",
				"name":       "",
				"symbol":     "",
				"sensibleid": "",
			}
			continue
		} else if err != nil {
			logger.Log.Info("getFTDecimal redis failed", zap.Error(err))
		}
		decimal, _ := strconv.Atoi(ftinfo["decimal"])
		ft.Decimal = decimal
		ft.Name = ftinfo["name"]
		ft.Symbol = ftinfo["symbol"]
		ft.SensibleIdHex = hex.EncodeToString([]byte(ftinfo["sensibleid"]))
	}
}

func ListFTInfoByGenesis(codeHashHex, genesisHex string) (ftRsp *model.FTInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT codehash, genesis, count(1),
       sum(in_data_value) AS in_volume , sum(out_data_value) AS out_volume,
       sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
WHERE code_type = 1 AND codehash = unhex('%s') AND genesis = unhex('%s')
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`, codeHashHex, genesisHex)
	ftsRsp, err := GetFTInfoBySQL(psql)
	if err != nil {
		return
	}
	if len(ftsRsp) > 0 {
		getFTDecimal(ftsRsp)
		return ftsRsp[0], nil
	}
	return nil, errors.New("not exist")
}

func GetFTSummary(codeHashHex string) (ftsRsp []*model.FTInfoResp, err error) {
	psql := fmt.Sprintf(`
SELECT codehash, genesis, count(1),
       sum(in_data_value) AS in_volume , sum(out_data_value) AS out_volume,
       sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
WHERE code_type = 1 AND codehash = unhex('%s')
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`, codeHashHex)
	ftsRsp, err = GetFTInfoBySQL(psql)
	if err != nil {
		return
	}
	getFTDecimal(ftsRsp)
	return
}

func GetFTInfo() (ftsRsp []*model.FTInfoResp, err error) {
	psql := `
SELECT codehash, genesis, count(1),
       sum(in_data_value) AS in_volume , sum(out_data_value) AS out_volume,
       sum(invalue) AS in_satoshi , sum(outvalue) AS out_satoshi FROM blk_codehash_height
WHERE code_type = 1
GROUP BY codehash, genesis
ORDER BY count(1) DESC
`
	ftsRsp, err = GetFTInfoBySQL(psql)
	if err != nil {
		return
	}
	getFTDecimal(ftsRsp)
	return
}

func GetFTInfoBySQL(psql string) (blksRsp []*model.FTInfoResp, err error) {
	blksRet, err := clickhouse.ScanAll(psql, ftInfoResultSRF)
	if err != nil {
		logger.Log.Info("query blk failed", zap.Error(err))
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
