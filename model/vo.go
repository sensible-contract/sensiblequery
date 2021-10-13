package model

////////////////
type Welcome struct {
	Contact string `json:"contact"`
	Job     string `json:"job"`
	Github  string `json:"github"`
}

////////////////

type MempoolInfoResp struct {
	TxCount int `json:"ntx"` // Mempool内包含的Tx数量
}

type BlockchainInfoResp struct {
	Chain         string `json:"chain"`         // main/test
	Blocks        int    `json:"blocks"`        // 最新区块总数
	Headers       int    `json:"headers"`       // 最新区块头总数
	BestBlockHash string `json:"bestBlockHash"` // 最新blockId
	Difficulty    string `json:"difficulty"`
	MedianTime    int    `json:"medianTime"`
	Chainwork     string `json:"chainwork"`
}

type BlockTokenVolumeResp struct {
	Height       int    `json:"height"` // 区块高度
	CodeHashHex  string `json:"codehash"`
	GenesisHex   string `json:"genesis"`
	CodeType     int    `json:"codeType"`     // 合约类型 0: None, 1: FT, 2: Unique, 3: NFT
	NFTIdx       int    `json:"nftIdx"`       // nft tokenIdx
	InDataValue  int    `json:"inDataValue"`  // 输入数量
	OutDataValue int    `json:"outDataValue"` // 输出数量
	InSatoshi    int    `json:"invalue"`
	OutSatoshi   int    `json:"outvalue"`
	BlockIdHex   string `json:"blkid"`
}

type BlockInfoResp struct {
	Height         int    `json:"height"`      // 当前区块高度
	BlockIdHex     string `json:"id"`          // 当前区块ID
	PrevBlockIdHex string `json:"prev"`        // 前一个区块ID
	NextBlockIdHex string `json:"next"`        // 下一个区块ID
	MerkleRootHex  string `json:"merkle"`      // Merkle Tree
	TxCount        int    `json:"ntx"`         // 区块内包含的Tx数量
	InSatoshi      int    `json:"inSatoshi"`   // 区块内输入额度总和，不包括区块reward
	OutSatoshi     int    `json:"outSatoshi"`  // 区块内输出额度总和，不包括区块reward/fee
	CoinbaseOut    int    `json:"coinbaseOut"` // 区块reward
	BlockTime      int    `json:"timestamp"`   // 区块时间戳
	Bits           int    `json:"bits"`
	BlockSize      int    `json:"size"` // 区块字节数
}

type TxInfoResp struct {
	TxIdHex    string `json:"txid"`
	InCount    int    `json:"nIn"`        // Tx包括的输入条数
	OutCount   int    `json:"nOut"`       // Tx包括的输出条数
	TxSize     int    `json:"size"`       // Tx字节数
	LockTime   int    `json:"locktime"`   // Tx Locktime
	InSatoshi  int    `json:"inSatoshi"`  // Tx输入的satoshi总和
	OutSatoshi int    `json:"outSatoshi"` // Tx输出的satoshi总和
	BlockTime  int    `json:"timestamp"`  // Tx所在区块的时间戳
	Height     int    `json:"height"`     // Tx所在区块的高度
	BlockIdHex string `json:"blkid"`      // Tx所在区块的Id
	Idx        int    `json:"idx"`        // Tx在区块中的顺序号
}

type TxInSpentResp struct {
	Height   int    `json:"height"` // 输出被花费的区块高度
	TxIdHex  string `json:"txid"`   // 输出被花费的txid
	Idx      int    `json:"idx"`    // 输出被花费的txid所在区块内序号
	UtxIdHex string `json:"utxid"`  // 输出txid参数
	Vout     int    `json:"vout"`   // 输出index参数
}

type TxInResp struct {
	Height       int    `json:"height"`    // Tx所在区块的高度
	TxIdHex      string `json:"txid"`      // Txid
	Idx          int    `json:"idx"`       // 输入index
	ScriptSigHex string `json:"scriptSig"` // scriptSig，Hex编码
	Sequence     int    `json:"sequence"`  // Tx input的sequence

	HeightTxo       int    `json:"heightTxo"`       // 当前输入花费的utxo所属的区块高度，如果为0则未花费
	UtxIdHex        string `json:"utxid"`           // 当前输入花费的outpoint的txid
	Vout            int    `json:"vout"`            // 当前输入花费的outpoint的index
	Address         string `json:"address"`         // 当前输入花费的outpoint的address
	IsNFT           bool   `json:"isNFT"`           // 当前输入花费是否为NFT
	CodeType        int    `json:"codeType"`        // 当前输出的合约类型 0: None, 1: FT, 2: Unique, 3: NFT
	CodeHashHex     string `json:"codehash"`        // 合约hash160(CodePart)
	GenesisHex      string `json:"genesis"`         // 当前输入花费的outpoint的genesis，Hex编码
	SensibleIdHex   string `json:"sensibleId"`      // 合约的sensibleId，即genesisTx的outpoint，Hex编码
	TokenId         string `json:"tokenId"`         // 当前输入的ft tokenId
	TokenAmount     string `json:"tokenAmount"`     // 当前输入花费的outpoint的ft tokenAmount
	TokenDecimal    int    `json:"tokenDecimal"`    // 当前输入花费的outpoint的ft decimal
	TokenName       string `json:"tokenName"`       // 当前输入的ft tokenName
	TokenSymbol     string `json:"tokenSymbol"`     // 当前输入的ft tokenSymbol
	TokenIndex      string `json:"tokenIndex"`      // 当前输入的nft tokenIndex
	MetaTxIdHex     string `json:"metaTxId"`        // 当前输入的nft metaTxId
	MetaOutputIndex int    `json:"metaOutputIndex"` // 当前输入的nft metaOutputIndex
	Satoshi         int    `json:"satoshi"`         // 当前输入花费的outpoint的satoshi
	ScriptTypeHex   string `json:"scriptType"`      // 当前输入锁定脚本类型，Hex编码
	ScriptPkHex     string `json:"scriptPk"`        // 当前输入锁定脚本，Hex编码
}

type TxOutResp struct {
	TxIdHex         string `json:"txid"`            // 当前txid
	Vout            int    `json:"vout"`            // 当前输出序号
	Address         string `json:"address"`         // 当前输出的address
	IsNFT           bool   `json:"isNFT"`           // 当前输出是否为NFT
	CodeType        int    `json:"codeType"`        // 当前输出的合约类型 0: None, 1: FT, 2: Unique, 3: NFT
	CodeHashHex     string `json:"codehash"`        // 合约hash160(CodePart)
	GenesisHex      string `json:"genesis"`         // 当前输出的genesis
	SensibleIdHex   string `json:"sensibleId"`      // 合约的sensibleId，即genesisTx的outpoint，Hex编码
	TokenId         string `json:"tokenId"`         // 当前输出的ft tokenId
	TokenAmount     string `json:"tokenAmount"`     // 当前输出的ft tokenAmount
	TokenDecimal    int    `json:"tokenDecimal"`    // 当前输出花费的outpoint的ft decimal
	TokenName       string `json:"tokenName"`       // 当前输出的ft tokenName
	TokenSymbol     string `json:"tokenSymbol"`     // 当前输出的ft tokenSymbol
	TokenIndex      string `json:"tokenIndex"`      // 当前输出的nft tokenIndex
	MetaTxIdHex     string `json:"metaTxId"`        // 当前输出的nft metaTxId
	MetaOutputIndex int    `json:"metaOutputIndex"` // 当前输出的nft metaOutputIndex
	Satoshi         int    `json:"satoshi"`         // 当前输出的satoshi
	ScriptTypeHex   string `json:"scriptType"`      // 当前输出锁定脚本类型
	ScriptPkHex     string `json:"scriptPk"`        // 当前输出锁定脚本
	Height          int    `json:"height"`          // 当前交易被打包的区块高度
	Idx             int    `json:"idx"`             // 当前交易被打包的区块内序号
}

type TxStandardOutResp struct {
	TxIdHex       string `json:"txid"`       // 当前txid
	Vout          int    `json:"vout"`       // 当前输出序号
	Satoshi       int    `json:"satoshi"`    // 当前输出的satoshi
	ScriptTypeHex string `json:"scriptType"` // 当前输出锁定脚本类型
	ScriptPkHex   string `json:"scriptPk"`   // 当前输出锁定脚本
	Height        int    `json:"height"`     // 当前交易被打包的区块高度
	Idx           int    `json:"idx"`        // 输出被花费的txid所在区块内序号
}

type AddressUTXOResp struct {
	Cursor                int                  `json:"cursor"`                // utxo结果偏移
	Total                 int                  `json:"total"`                 // utxo总量 total = confirmed + unconfirmed - unconfirmedSpend
	TotalConfirmed        int                  `json:"totalConfirmed"`        // 已确认utxo总量
	TotalUnconfirmedNew   int                  `json:"totalUnconfirmed"`      // 未确认新utxo总量
	TotalUnconfirmedSpend int                  `json:"totalUnconfirmedSpend"` // 未确认已花费utxo总量
	UTXO                  []*TxStandardOutResp `json:"utxo"`                  // utxo结果
}

type AddressTokenUTXOResp struct {
	Cursor                int          `json:"cursor"`                // utxo结果偏移
	Total                 int          `json:"total"`                 // utxo总量 total = confirmed + unconfirmed - unconfirmedSpend
	TotalConfirmed        int          `json:"totalConfirmed"`        // 已确认utxo总量
	TotalUnconfirmedNew   int          `json:"totalUnconfirmed"`      // 未确认新utxo总量
	TotalUnconfirmedSpend int          `json:"totalUnconfirmedSpend"` // 未确认已花费utxo总量
	UTXO                  []*TxOutResp `json:"utxo"`                  // utxo结果
}

type TxOutHistoryResp struct {
	TxOutResp
	BlockTime int `json:"timestamp"` // 区块时间戳
	IOType    int `json:"ioType"`    // 1为输出包含(即收入)，0为输入包含(即花费)
}

type TxOutStatusResp struct {
	TxOutResp

	TxIdSpentHex string `json:"txidSpent"`   // 当前输出被花费的txid
	HeightSpent  int    `json:"heightSpent"` // 当前输出被花费的区块高度，如果为0则未花费
}

type BalanceResp struct {
	Address        string `json:"address"`        // address
	Satoshi        int    `json:"satoshi"`        // 余额satoshi
	PendingSatoshi int    `json:"pendingSatoshi"` // 待确认余额satoshi
	UtxoCount      int    `json:"utxoCount"`      // UTXO 数量
}

////////////////
type ContractSwapDataResp struct {
	Height          int    `json:"height"`    // 区块高度
	BlockTime       int    `json:"timestamp"` // 区块时间戳
	CodeType        int    `json:"codeType"`  // 合约类型 0: None, 1: FT, 2: Unique, 3: NFT
	Operation       int    `json:"operation"` // 0: sell, 1: buy, 2: add, 3: remove
	InToken1Amount  int    `json:"inToken1Amount"`
	InToken2Amount  int    `json:"inToken2Amount"`
	InLpAmount      int    `json:"inTokenLpAmount"`
	OutToken1Amount int    `json:"outToken1Amount"`
	OutToken2Amount int    `json:"outToken2Amount"`
	OutLpAmount     int    `json:"outTokenLpAmount"`
	Idx             int    `json:"idx"`
	TxIdHex         string `json:"txid"` // 当前txid
}

type ContractSwapAggregateResp struct {
	Height       int     `json:"height"`       // 区块高度
	BlockTime    int     `json:"timestamp"`    // 区块时间戳
	OpenPrice    float64 `json:"openPrice"`    // 开盘价格
	ClosePrice   float64 `json:"closePrice"`   // 收盘价格
	MinPrice     float64 `json:"minPrice"`     // 最低价格
	MaxPrice     float64 `json:"maxPrice"`     // 最高价格
	Token1Volume int     `json:"token1Volume"` // token1交易量
	Token2Volume int     `json:"token2Volume"` // token2交易量
}

type ContractSwapAggregateAmountResp struct {
	Height      int `json:"height"`      // 区块高度
	BlockTime   int `json:"timestamp"`   // 区块时间戳
	OpenAmount  int `json:"openAmount"`  // 开盘Token1存量
	CloseAmount int `json:"closeAmount"` // 收盘Token1存量
	MinAmount   int `json:"minAmount"`   // 最低Token1存量
	MaxAmount   int `json:"maxAmount"`   // 最高Token1存量
	Count       int `json:"txCount"`     // 交易次数
}
