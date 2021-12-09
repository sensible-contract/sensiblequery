package service

import (
	"encoding/hex"
	"sensiblequery/lib/blkparser"
	"sensiblequery/lib/utils"
	"sensiblequery/logger"
	"sensiblequery/model"

	redis "github.com/go-redis/redis/v8"
	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	"go.uber.org/zap"
)

//////////////// address utxo
func GetNFTAuctionUtxoByNFTIDMerge(codeHash, nftId []byte, isReadyOnly bool) (nftAuctionsRsp []*model.NFTAuctionResp, err error) {
	// fixme: 可能被恶意创建sell utxo
	key := "mp:{nau" + string(nftId) + string(codeHash) + "}"
	respMempool, err := GetNFTAuctionUtxoByKey(key)
	if err != nil {
		return nil, err
	}

	key = "{nau" + string(nftId) + string(codeHash) + "}"
	resp, err := GetNFTAuctionUtxoByKey(key)
	if err != nil {
		return nil, err
	}

	nftAuctionsRsp = make([]*model.NFTAuctionResp, 0)
	for _, data := range respMempool {
		if isReadyOnly && !data.IsReady {
			continue
		}
		nftAuctionsRsp = append(nftAuctionsRsp, data)
	}

	for _, data := range resp {
		if isReadyOnly && !data.IsReady {
			continue
		}
		nftAuctionsRsp = append(nftAuctionsRsp, data)
	}

	return
}

func GetNFTAuctionUtxoByKey(key string) (nftAuctionsRsp []*model.NFTAuctionResp, err error) {
	utxoOutpoints, err := rdb.ZRange(ctx, key, 0, 16).Result()
	if err != nil {
		logger.Log.Info("GetNFTAuctionUtxoByKey redis failed", zap.Error(err))
		return
	}
	nftAuctionsRsp, err = getNFTAuctionUtxoFromRedis(utxoOutpoints)
	if err != nil {
		return nil, err
	}
	return nftAuctionsRsp, nil
}

////////////////
func getNFTAuctionUtxoFromRedis(utxoOutpoints []string) (nftAuctionsRsp []*model.NFTAuctionResp, err error) {
	logger.Log.Info("getNFTAuctionUtxoFromRedis redis", zap.Int("nUTXO", len(utxoOutpoints)))
	nftAuctionsRsp = make([]*model.NFTAuctionResp, 0)
	pipe := rdb.Pipeline()

	outpointsCmd := make([]*redis.StringCmd, 0)
	for _, outpoint := range utxoOutpoints {
		outpointsCmd = append(outpointsCmd, pipe.Get(ctx, "u"+outpoint))
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for outpointIdx, data := range outpointsCmd {
		outpoint := utxoOutpoints[outpointIdx]
		res, err := data.Result()
		if err == redis.Nil {
			logger.Log.Info("redis not found", zap.String("outpoint", hex.EncodeToString([]byte(outpoint))))
			continue
		} else if err != nil {
			panic(err)
		}

		txout := model.NewTxoData([]byte(outpoint), []byte(res))
		nftAuctionRsp := &model.NFTAuctionResp{
			TxIdHex: blkparser.HashString(txout.UTxid),
			Vout:    int(txout.Vout),
			Satoshi: int(txout.Satoshi),
			Height:  int(txout.BlockHeight),
			Idx:     int(txout.TxIdx),
		}
		txo := scriptDecoder.ExtractPkScriptForTxo(txout.Script, txout.ScriptType)
		if txo.CodeType != scriptDecoder.CodeType_NONE && txo.CodeType != scriptDecoder.CodeType_SENSIBLE {
			nftAuctionRsp.CodeHashHex = hex.EncodeToString(txo.CodeHash[:])
		}
		if txo.CodeType == scriptDecoder.CodeType_NFT_AUCTION {
			nftAuctionRsp.SensibleIdHex = hex.EncodeToString(txo.NFTAuction.SensibleId[:])
			nftAuctionRsp.GenesisHex = nftAuctionRsp.SensibleIdHex

			nftAuctionRsp.NFTCodeHashHex = hex.EncodeToString(txo.NFTAuction.NFTCodeHash[:])
			nftAuctionRsp.NFTIDHex = hex.EncodeToString(txo.NFTAuction.NFTID[:])
			nftAuctionRsp.FeeAmount = int(txo.NFTAuction.FeeAmount)
			nftAuctionRsp.FeeAddress = utils.EncodeAddress(txo.NFTAuction.FeeAddressPkh[:], utils.PubKeyHashAddrID)
			nftAuctionRsp.StartBsvPrice = int(txo.NFTAuction.StartBsvPrice)
			nftAuctionRsp.SenderAddress = utils.EncodeAddress(txo.NFTAuction.SenderAddressPkh[:], utils.PubKeyHashAddrID)
			nftAuctionRsp.EndTimestamp = int(txo.NFTAuction.EndTimestamp)
			nftAuctionRsp.BidTimestamp = int(txo.NFTAuction.BidTimestamp)
			nftAuctionRsp.BidBsvPrice = int(txo.NFTAuction.BidBsvPrice)
			nftAuctionRsp.BidderAddress = utils.EncodeAddress(txo.NFTAuction.BidderAddressPkh[:], utils.PubKeyHashAddrID)

			// 设置准备状态
			contractHashAsAddressPkh := blkparser.GetHash160(txout.Script)
			countRsp, err := GetNFTCountByCodeHashGenesisAddress(
				txo.NFTAuction.NFTCodeHash[:],
				txo.NFTAuction.NFTID[:], contractHashAsAddressPkh)
			if err == nil && countRsp.Count+countRsp.PendingCount > 0 {
				nftAuctionRsp.IsReady = true
			}
		}

		nftAuctionsRsp = append(nftAuctionsRsp, nftAuctionRsp)
	}

	return nftAuctionsRsp, nil
}
