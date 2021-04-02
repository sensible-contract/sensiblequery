package model

type TokenCodeHashResp struct {
	CodeHashHex string `json:"codehash"` // FT合约hash160(CodePart)
	Count       int    `json:"count"`    // 采用当前合约的Token种类数，NFT包括具体NFT数量；FT只包括种类数量
	InTimes     int    `json:"inTimes"`  // 总输入次数
	OutTimes    int    `json:"outTimes"` // 总输出次数
}
