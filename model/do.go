package model

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
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	ScriptPk   []byte `db:"script_pk"`
}

type TxOutDO struct {
	TxId       []byte `db:"txid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	ScriptPk   []byte `db:"script_pk"`
	Height     uint32 `db:"height"`
}

type TxOutHistoryDO struct {
	TxId       []byte `db:"txid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	Height     uint32 `db:"height"`
	Idx        uint32 `db:"txidx"`
	IOType     uint8  `db:"io_type"` // 0: input; 1: output
}

type TxOutStatusDO struct {
	TxId       []byte `db:"txid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	ScriptPk   []byte `db:"script_pk"`
	Height     uint32 `db:"height"`

	TxIdSpent   []byte `db:"txid_spent"`
	HeightSpent uint32 `db:"height_spent"`
}
