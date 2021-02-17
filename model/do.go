package model

////////////////
type BlockDO struct {
	Height      uint32 `db:"height"`
	BlockId     []byte `db:"blkid"`
	PrevBlockId []byte `db:"previd"`
	TxCount     uint64 `db:"ntx"`
}

type TxDO struct {
	TxId     []byte `db:"txid"`
	InCount  uint32 `db:"nin"`
	OutCount uint32 `db:"nout"`
	Height   uint32 `db:"height"`
	BlockId  []byte `db:"blkid"`
	Idx      uint64 `db:"idx"`
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
	ScriptSig []byte `db:"script_sig"`

	HeightTxo  uint32 `db:"height_txo"`
	UtxId      []byte `db:"utxid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
}

type TxOutDO struct {
	TxId       []byte `db:"txid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	Script     []byte `db:"script"`
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
	IOType     uint8  `db:"io_type"` // 0: input; 1: output
}

type TxOutStatusDO struct {
	TxId       []byte `db:"txid"`
	Vout       uint32 `db:"vout"`
	Address    []byte `db:"address"`
	Genesis    []byte `db:"genesis"`
	Satoshi    uint64 `db:"satoshi"`
	ScriptType []byte `db:"script_type"`
	Script     []byte `db:"script"`
	Height     uint32 `db:"height"`

	TxIdSpent   []byte `db:"txid_spent"`
	HeightSpent uint32 `db:"height_spent"`
}
