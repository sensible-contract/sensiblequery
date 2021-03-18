package model

type FTInfoResp struct {
	CodeHashHex string `json:"codehash"` // FT合约hash160(CodePart)
	GenesisHex  string `json:"genesis"`  // FT合约的genesis，Hex编码
	Name        string `json:"name"`     // FT name
	Symbol      string `json:"symbol"`   // FT symbol
	Decimal     int    `json:"decimal"`  // decimal
	Desc        string `json:"desc"`     // FT 描述
	Icon        string `json:"icon"`     // FT icon url
	Website     string `json:"website"`  // FT website url
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

type NFTOwnerByAddressResp struct {
	CodeHashHex string `json:"codehash"` // NFT合约hash160(CodePart)
	GenesisHex  string `json:"genesis"`  // NFT合约的genesis，Hex编码
	Symbol      string `json:"symbol"`   // NFT symbol
	Count       int    `json:"count"`    // 持有的当前NFT个数
}
