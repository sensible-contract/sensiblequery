package model

type TokenInfoResp struct {
	CodeCount int `json:"countCode"` // 当前合约代码种类数

	NFTCount      int `json:"countNFT"`      // 当前NFT总数
	NFTIDCount    int `json:"countNFTID"`    // 当前NFTID总数
	InTimesNFT    int `json:"inTimesNFT"`    // NFT总输入次数
	OutTimesNFT   int `json:"outTimesNFT"`   // NFT总输出次数
	OwnerNFTCount int `json:"countOwnerNFT"` // 当前持有NFT人数

	FTCount      int `json:"countFT"`      // 当前FT总数
	InTimesFT    int `json:"inTimesFT"`    // FT总输入次数
	OutTimesFT   int `json:"outTimesFT"`   // FT总输出次数
	OwnerFTCount int `json:"countOwnerFT"` // 当前持有FT人数
}

type TokenCodeHashResp struct {
	CodeHashHex string `json:"codehash"` // FT合约hash160(CodePart)
	Count       int    `json:"count"`    // 采用当前合约的Token种类数，NFT包括具体NFT数量；FT只包括种类数量
	InTimes     int    `json:"inTimes"`  // 总输入次数
	OutTimes    int    `json:"outTimes"` // 总输出次数
}
