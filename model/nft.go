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

type NFTSellResp struct {
	Height          int    `json:"height"`          // 当前交易被打包的区块高度
	Idx             int    `json:"idx"`             // 输出被花费的txid所在区块内序号
	TxIdHex         string `json:"txid"`            // 售卖合约txid
	Vout            int    `json:"vout"`            // 售卖合约输出序号
	Satoshi         int    `json:"satoshi"`         // 售卖合约输出的satoshi
	Address         string `json:"address"`         // 当前售卖人seller的address
	CodeHashHex     string `json:"codehash"`        // 当前售卖NFT合约hash160(CodePart)
	GenesisHex      string `json:"genesis"`         // 当前售卖NFT的genesis
	SensibleIdHex   string `json:"sensibleId"`      // 当前售卖NFT合约的sensibleId，即genesisTx的outpoint，Hex编码
	TokenIndex      string `json:"tokenIndex"`      // 当前售卖NFT的tokenIndex
	MetaTxIdHex     string `json:"metaTxId"`        // 当前售卖NFT的metaTxId
	MetaOutputIndex int    `json:"metaOutputIndex"` // 当前售卖NFT的metaOutputIndex
	Price           int    `json:"price"`           // 当前售卖NFT的出售价格(satoshi)
	IsReady         bool   `json:"isReady"`         // 当前售卖NFT是否已准备好(转出到售卖合约)
}
