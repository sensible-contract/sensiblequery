package model

type FTInfoResp struct {
	CodeHashHex string `json:"codehash"`  // FT合约hash160(CodePart)
	GenesisHex  string `json:"genesis"`   // FT合约的genesis，Hex编码
	Name        string `json:"name"`      // FT name
	Symbol      string `json:"symbol"`    // FT symbol
	Decimal     int    `json:"decimal"`   // decimal
	Desc        string `json:"desc"`      // FT 描述
	Icon        string `json:"icon"`      // FT icon url
	Website     string `json:"website"`   // FT website url
	Count       int    `json:"count"`     // 出现此ft合约的区块次数
	InVolume    int    `json:"inVolume"`  // 输入数量
	OutVolume   int    `json:"outVolume"` // 输出数量
	InSatoshi   int    `json:"inSatoshi"`
	OutSatoshi  int    `json:"outSatoshi"`
}

type FTTransferVolumeResp struct {
	Height    int `json:"height"`    // 区块高度
	InVolume  int `json:"inVolume"`  // 输入数量
	OutVolume int `json:"outVolume"` // 输出数量
}

type FTOwnerBalanceResp struct {
	Address        string `json:"address"`        // token持有人的address
	Balance        int    `json:"balance"`        // 余额
	PendingBalance int    `json:"pendingBalance"` // 待确认余额
	Decimal        int    `json:"decimal"`        // decimal
}

type FTOwnerBalanceWithUtxoCountResp struct {
	Address        string `json:"address"`        // token持有人的address
	Balance        int    `json:"balance"`        // 余额
	PendingBalance int    `json:"pendingBalance"` // 待确认余额
	Decimal        int    `json:"decimal"`        // decimal
	UtxoCount      int    `json:"utxoCount"`      // UTXO 数量
}

type FTSummaryByAddressResp struct {
	CodeHashHex    string `json:"codehash"`       // FT合约hash160(CodePart)
	GenesisHex     string `json:"genesis"`        // FT合约的genesis，Hex编码
	Name           string `json:"name"`           // FT name
	Symbol         string `json:"symbol"`         // FT symbol
	Balance        int    `json:"balance"`        // 余额
	PendingBalance int    `json:"pendingBalance"` // 待确认余额
	Decimal        int    `json:"decimal"`        // decimal
}
