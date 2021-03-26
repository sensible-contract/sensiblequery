package model

type NFTInfoResp struct {
	CodeHashHex string `json:"codehash"` // NFT合约hash160(CodePart)
	GenesisHex  string `json:"genesis"`  // NFT合约的genesis，Hex编码
	Name        string `json:"name"`     // NFT name
	Symbol      string `json:"symbol"`   // NFT symbol
	Desc        string `json:"desc"`     // NFT 描述
	Icon        string `json:"icon"`     // NFT icon url
	Website     string `json:"website"`  // NFT website url
	Count       int    `json:"count"`    // 当前NFT个数
	InTimes     int    `json:"inTimes"`  // 总输入次数
	OutTimes    int    `json:"outTimes"` // 总输出次数
	InSatoshi   int    `json:"inSatoshi"`
	OutSatoshi  int    `json:"outSatoshi"`
}

type NFTTransferTimesResp struct {
	Height   int `json:"height"`   // 区块高度
	InTimes  int `json:"inTimes"`  // 输入次数
	OutTimes int `json:"outTimes"` // 输出次数
}

type NFTOwnerResp struct {
	Address string `json:"address"` // token持有人的address
	TokenId int    `json:"tokenId"` // 持有的NFT id
}

type NFTSummaryResp struct {
	Address string `json:"address"` // token持有人的address
	Count   int    `json:"count"`   // 持有的当前NFT个数
}

type NFTOwnerByAddressResp struct {
	CodeHashHex string `json:"codehash"` // NFT合约hash160(CodePart)
	GenesisHex  string `json:"genesis"`  // NFT合约的genesis，Hex编码
	Symbol      string `json:"symbol"`   // NFT symbol
	Count       int    `json:"count"`    // 持有的当前NFT个数
}
