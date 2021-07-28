package model

type NFTInfoResp struct {
	CodeHashHex     string `json:"codehash"`        // NFT合约hash160(CodePart)
	GenesisHex      string `json:"genesis"`         // NFT合约的genesis，Hex编码
	SensibleIdHex   string `json:"sensibleId"`      // NFT合约的sensibleId，即genesisTx的outpoint，Hex编码
	MetaTxIdHex     string `json:"metaTxId"`        // NFT合约的 MetaTxId
	MetaOutputIndex int    `json:"metaOutputIndex"` // NFT合约的 MetaOutputIndex
	Name            string `json:"name"`            // NFT name
	Symbol          string `json:"symbol"`          // NFT symbol
	Desc            string `json:"desc"`            // NFT 描述
	Icon            string `json:"icon"`            // NFT icon url
	Website         string `json:"website"`         // NFT website url
	Supply          int    `json:"supply"`          // 当前NFT最大发行量
	Count           int    `json:"count"`           // 当前NFT个数
	InTimes         int    `json:"inTimes"`         // 总输入次数
	OutTimes        int    `json:"outTimes"`        // 总输出次数
	InSatoshi       int    `json:"inSatoshi"`
	OutSatoshi      int    `json:"outSatoshi"`
}

type NFTTransferTimesResp struct {
	Height   int `json:"height"`   // 区块高度
	InTimes  int `json:"inTimes"`  // 输入次数
	OutTimes int `json:"outTimes"` // 输出次数
}

type NFTOwnerResp struct {
	Address      string `json:"address"`      // token持有人的address
	Count        int    `json:"count"`        // 持有的当前NFT个数
	PendingCount int    `json:"pendingCount"` // 待确认的当前NFT个数
}

type NFTSummaryByAddressResp struct {
	CodeHashHex     string `json:"codehash"`        // NFT合约hash160(CodePart)
	GenesisHex      string `json:"genesis"`         // NFT合约的genesis，Hex编码
	SensibleIdHex   string `json:"sensibleId"`      // NFT合约的sensibleId，即genesisTx的outpoint，Hex编码
	MetaTxIdHex     string `json:"metaTxId"`        // NFT合约的 MetaTxId
	MetaOutputIndex int    `json:"metaOutputIndex"` // NFT合约的 MetaOutputIndex
	Symbol          string `json:"symbol"`          // NFT symbol
	Supply          int    `json:"supply"`          // 当前NFT最大发行量
	Count           int    `json:"count"`           // 持有的当前NFT个数
	PendingCount    int    `json:"pendingCount"`    // 待确认的当前NFT个数
}
