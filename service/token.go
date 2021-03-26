package service

import (
	"encoding/hex"
	"log"
	"satoblock/lib/utils"
	"satoblock/model"
)

func GetTokenOwnersByCodeHashGenesis(cursor, size int, codeHash, genesisId []byte) (ftOwnersRsp []*model.FTOwnerBalanceResp, err error) {
	vals, err := rdb.ZRevRangeWithScores("fs"+string(codeHash)+string(genesisId), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ftOwnersRsp = append(ftOwnersRsp, &model.FTOwnerBalanceResp{
			Address: utils.EncodeAddress(val.Member.([]byte), utils.PubKeyHashAddrIDMainNet),
			Balance: int(val.Score),
		})
	}

	return ftOwnersRsp, nil
}

func GetAllTokenBalanceByAddress(cursor, size int, addressPkh []byte) (ftOwnersRsp []*model.FTOwnerByAddressResp, err error) {
	vals, err := rdb.ZRevRangeWithScores("fs"+string(addressPkh), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ftOwnersRsp = append(ftOwnersRsp, &model.FTOwnerByAddressResp{
			CodeHashHex: hex.EncodeToString(val.Member.([]byte)[:20]),
			GenesisHex:  hex.EncodeToString(val.Member.([]byte)[20:]),
			Balance:     int(val.Score),
		})
	}

	return ftOwnersRsp, nil
}

func GetTokenBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte) (balanceRsp *model.FTOwnerBalanceResp, err error) {
	score, err := rdb.ZScore("fb"+string(codeHash)+string(genesisId), string(addressPkh)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	balanceRsp = &model.FTOwnerBalanceResp{
		Address: utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrIDMainNet),
		Balance: int(score),
	}
	return balanceRsp, nil
}

////////////////
// nft

func GetNFTOwnersByCodeHashGenesis(cursor, size int, codeHash, genesisId []byte) (ownersRsp []*model.NFTSummaryResp, err error) {
	vals, err := rdb.ZRevRangeWithScores("ns"+string(codeHash)+string(genesisId), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ownersRsp = append(ownersRsp, &model.NFTSummaryResp{
			Address: utils.EncodeAddress(val.Member.([]byte), utils.PubKeyHashAddrIDMainNet),
			Count:   int(val.Score),
		})
	}

	return ownersRsp, nil
}

func GetAllNFTBalanceByAddress(cursor, size int, addressPkh []byte) (ftOwnersRsp []*model.FTOwnerByAddressResp, err error) {
	vals, err := rdb.ZRevRangeWithScores("ns"+string(addressPkh), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	for _, val := range vals {
		ftOwnersRsp = append(ftOwnersRsp, &model.FTOwnerByAddressResp{
			CodeHashHex: hex.EncodeToString(val.Member.([]byte)[:20]),
			GenesisHex:  hex.EncodeToString(val.Member.([]byte)[20:]),
			Balance:     int(val.Score),
		})
	}

	return ftOwnersRsp, nil
}

func GetNFTBalanceByCodeHashGenesisAddress(codeHash, genesisId, addressPkh []byte) (balanceRsp *model.FTOwnerBalanceResp, err error) {
	score, err := rdb.ZScore("no"+string(codeHash)+string(genesisId), string(addressPkh)).Result()
	if err != nil {
		log.Printf("GetFTOwnersByCodeHashGenesis redis failed: %v", err)
		return
	}

	balanceRsp = &model.FTOwnerBalanceResp{
		Address: utils.EncodeAddress(addressPkh, utils.PubKeyHashAddrIDMainNet),
		Balance: int(score),
	}
	return balanceRsp, nil
}
