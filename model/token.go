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
