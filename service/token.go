package service

import (
	"encoding/hex"
	"log"
	"satosensible/lib/utils"
	"satosensible/model"
)

////////////////
// ft
func GetTokenOwnersByCodeHashGenesis(cursor, size int, codeHash, genesisId []byte) (ftOwnersRsp []*model.FTOwnerBalanceResp, err error) {
	vals, err := rdb.ZRevRangeWithScores(ctx, "fb"+string(codeHash)+string(genesisId), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ftOwnersRsp = append(ftOwnersRsp, &model.FTOwnerBalanceResp{
			Address: utils.EncodeAddress([]byte(val.Member.(string)), utils.PubKeyHashAddrID),
			Balance: int(val.Score),
		})
	}

	return ftOwnersRsp, nil
}

func GetAllTokenBalanceByAddress(cursor, size int, addressPkh []byte) (ftOwnersRsp []*model.FTOwnerByAddressResp, err error) {
	vals, err := rdb.ZRevRangeWithScores(ctx, "fs"+string(addressPkh), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ftOwnersRsp = append(ftOwnersRsp, &model.FTOwnerByAddressResp{
			CodeHashHex: hex.EncodeToString([]byte(val.Member.(string))[:20]),
			GenesisHex:  hex.EncodeToString([]byte(val.Member.(string))[20:]),
			Balance:     int(val.Score),
		})
	}

	return ftOwnersRsp, nil
}

func GetTokenBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte) (balanceRsp *model.FTOwnerBalanceResp, err error) {
	score, err := rdb.ZScore(ctx, "fb"+string(codeHash)+string(genesisId), string(addressPkh)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	balanceRsp = &model.FTOwnerBalanceResp{
		Address: utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrID),
		Balance: int(score),
	}
	return balanceRsp, nil
}

////////////////
// nft
func GetNFTOwnersByCodeHashGenesis(cursor, size int, codeHash, genesisId []byte) (ownersRsp []*model.NFTSummaryResp, err error) {
	vals, err := rdb.ZRevRangeWithScores(ctx, "no"+string(codeHash)+string(genesisId), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ownersRsp = append(ownersRsp, &model.NFTSummaryResp{
			Address: utils.EncodeAddress([]byte(val.Member.(string)), utils.PubKeyHashAddrID),
			Count:   int(val.Score),
		})
	}

	return ownersRsp, nil
}

func GetAllNFTBalanceByAddress(cursor, size int, addressPkh []byte) (ftOwnersRsp []*model.FTOwnerByAddressResp, err error) {
	vals, err := rdb.ZRevRangeWithScores(ctx, "ns"+string(addressPkh), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ftOwnersRsp = append(ftOwnersRsp, &model.FTOwnerByAddressResp{
			CodeHashHex: hex.EncodeToString([]byte(val.Member.(string))[:20]),
			GenesisHex:  hex.EncodeToString([]byte(val.Member.(string))[20:]),
			Balance:     int(val.Score),
		})
	}

	return ftOwnersRsp, nil
}

func GetNFTCountByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte) (balanceRsp *model.NFTSummaryResp, err error) {
	score, err := rdb.ZScore(ctx, "no"+string(codeHash)+string(genesisId), string(addressPkh)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	balanceRsp = &model.NFTSummaryResp{
		Address: utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrID),
		Count:   int(score),
	}
	return balanceRsp, nil
}
