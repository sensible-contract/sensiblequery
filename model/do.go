package model

type TokenInfoDO struct {
	CodeHash   []byte `db:"codehash"`  // Token合约hash160(CodePart)
	Genesis    []byte `db:"genesis"`   // Token合约的genesis，Hex编码
	Name       string `db:"name"`      // Token name
	Symbol     string `db:"symbol"`    // Token symbol
	Desc       string `db:"desc"`      // Token 描述
	Icon       string `db:"icon"`      // Token icon url
	Website    string `db:"website"`   // Token website url
	Count      uint64 `db:"count"`     // 当前Token个数
	InTimes    uint64 `db:"in_times"`  // 总输入次数
	OutTimes   uint64 `db:"out_times"` // 总输出次数
	InSatoshi  uint64 `db:"invalue"`
	OutSatoshi uint64 `db:"outvalue"`
}

type TokenCodeHashDO struct {
	CodeHash []byte `db:"codehash"`  // Token合约hash160(CodePart)
	Count    uint64 `db:"count"`     // 采用当前合约的Token种类数，NFT包括具体NFT数量；FT只包括种类数量
	InTimes  uint64 `db:"in_times"`  // 总输入次数
	OutTimes uint64 `db:"out_times"` // 总输出次数
}

type NFTInfoDO struct {
	CodeHash   []byte `db:"codehash"`  // NFT合约hash160(CodePart)
	Genesis    []byte `db:"genesis"`   // NFT合约的genesis，Hex编码
	Name       string `db:"name"`      // NFT name
	Symbol     string `db:"symbol"`    // NFT symbol
	Desc       string `db:"desc"`      // NFT 描述
	Icon       string `db:"icon"`      // NFT icon url
	Website    string `db:"website"`   // NFT website url
	Count      uint64 `db:"count"`     // 当前NFT个数
	InTimes    uint64 `db:"in_times"`  // 总输入次数
	OutTimes   uint64 `db:"out_times"` // 总输出次数
	InSatoshi  uint64 `db:"invalue"`
	OutSatoshi uint64 `db:"outvalue"`
}

type FTInfoDO struct {
	CodeHash   []byte `json:"codehash"`  // FT合约hash160(CodePart)
	Genesis    []byte `json:"genesis"`   // FT合约的genesis，Hex编码
	Name       string `json:"name"`      // FT name
	Symbol     string `json:"symbol"`    // FT symbol
	Decimal    int    `json:"decimal"`   // decimal
	Desc       string `json:"desc"`      // FT 描述
	Icon       string `json:"icon"`      // FT icon url
	Website    string `json:"website"`   // FT website url
	Count      uint64 `json:"count"`     // 出现此合约的区块次数
	InVolume   uint64 `json:"inVolume"`  // 输入数量
	OutVolume  uint64 `json:"outVolume"` // 输出数量
	InSatoshi  uint64 `json:"inSatoshi"`
	OutSatoshi uint64 `json:"outSatoshi"`
}

type BlockTokenVolumeDO struct {
	Height       uint32 `db:"height"` // 区块高度
	CodeHash     []byte `db:"codehash"`
	Genesis      []byte `db:"genesis"`
	CodeType     uint32 `db:"code_type"`      // 合约类型 0: None, 1: FT, 2: UNIQUE, 3: NFT
	NFTIdx       uint64 `db:"nft_idx"`        // nft tokenIdx
	InDataValue  uint64 `db:"in_data_value"`  // 输入数量
	OutDataValue uint64 `db:"out_data_value"` // 输出数量
	InSatoshi    uint64 `db:"invalue"`
	OutSatoshi   uint64 `db:"outvalue"`
	BlockId      []byte `db:"blkid"`
}

////////////////
type BlockDO struct {
	Height      uint32 `db:"height"`
	BlockId     []byte `db:"blkid"`
	PrevBlockId []byte `db:"previd"`
	NextBlockId []byte `db:"nextid"`
	MerkleRoot  []byte `db:"merkle"`
	TxCount     uint64 `db:"ntx"`
	InSatoshi   uint64 `db:"invalue"`
	OutSatoshi  uint64 `db:"outvalue"`
	CoinbaseOut uint64 `db:"coinbase_out"`
	BlockTime   uint32 `db:"blocktime"`
	Bits        uint32 `db:"bits"`
	BlockSize   uint32 `db:"blocksize"`
}

type TxDO struct {
	TxId       []byte `db:"txid"`
	InCount    uint32 `db:"nin"`
	OutCount   uint32 `db:"nout"`
	TxSize     uint32 `db:"txsize"`
	LockTime   uint32 `db:"locktime"`
	InSatoshi  uint64 `db:"invalue"`
	OutSatoshi uint64 `db:"outvalue"`
	BlockTime  uint32 `db:"blocktime"`
	Height     uint32 `db:"height"`
	BlockId    []byte `db:"blkid"`
	Idx        uint64 `db:"idx"`
}

type TxInSpentDO struct {
	Height uint32 `db:"height"`
	TxId   []byte `db:"txid"`
	Idx    uint32 `db:"idx"`
	UtxId  []byte `db:"utxid"`
	Vout   uint32 `db:"vout"`
}

type TxInDO struct {
	Height    uint32 `db:"height"`
	TxId      []byte `db:"txid"`
	Idx       uint32 `db:"idx"`
	ScriptSig []byte `db:"scriptSig"`
	Sequence  uint32 `db:"nsequence"`

	HeightTxo  uint32 `db:"height_txo"`
	UtxId      []byte `db:"utxid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	CodeHash   []byte `db:"codehash"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	ScriptPk   []byte `db:"script_pk"`
}

type TxOutDO struct {
	TxId       []byte `db:"txid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	CodeHash   []byte `db:"codehash"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	ScriptPk   []byte `db:"script_pk"`
	Height     uint32 `db:"height"`
	Idx        uint32 `db:"txidx"`
}

type TxOutHistoryDO struct {
	TxOutDO

	IOType uint8 `db:"io_type"` // 0: input; 1: output
}

type TxOutStatusDO struct {
	TxOutDO

	TxIdSpent   []byte `db:"txid_spent"`
	HeightSpent uint32 `db:"height_spent"`
}
