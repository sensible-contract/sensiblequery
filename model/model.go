package model

import "encoding/json"

type TxRequest struct {
	TxHex   string `json:"txHex"`
	ByTxHex string `json:"byTxHex"`
}

type TxResponse struct {
	TxId    string `json:"txId"`
	Index   int    `json:"index"`
	ByTxId  string `json:"byTxId"`
	Sig     string `json:"sigBE"`
	Padding string `json:"padding"`
	Payload string `json:"payload"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (t *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}

////////////////
type Welcome struct {
	PubKey  string `json:"pubKey"`
	Contact string `json:"contact"`
	Job     string `json:"job"`
	Github  string `json:"github"`
}

////////////////

type BlockchainInfoResp struct {
	Chain         string `json:"chain"`
	Blocks        int    `json:"blocks"`
	Headers       int    `json:"headers"`
	BestBlockHash string `json:"bestBlockHash"`
	Difficulty    string `json:"difficulty"`
	MedianTime    int    `json:"medianTime"`
	Chainwork     string `json:"chainwork"`
}

type BlockInfoResp struct {
	Height         int    `json:"height"`
	BlockIdHex     string `json:"id"`
	PrevBlockIdHex string `json:"prev"`
	TxCount        int    `json:"ntx"`
}

type TxInfoResp struct {
	TxIdHex  string `json:"txid"`
	InCount  int    `json:"nIn"`
	OutCount int    `json:"nOut"`

	Height     int    `json:"height"`
	BlockIdHex string `json:"blkid"`
	Idx        int    `json:"idx"`
}

type TxInSpentResp struct {
	Height   int    `json:"height"`
	TxIdHex  string `json:"txid"`
	Idx      int    `json:"idx"`
	UtxIdHex string `json:"utxid"`
	Vout     int    `json:"vout"`
}

type TxInResp struct {
	Height       int    `json:"height"`
	TxIdHex      string `json:"txid"`
	Idx          int    `json:"idx"`
	ScriptSigHex string `json:"script_sig"`

	HeightTxo     int    `json:"height_txo"`
	UtxIdHex      string `json:"utxid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	GenesisHex    string `json:"genesis"`
	Satoshi       int    `json:"satoshi"`
	ScriptTypeHex string `json:"script_type"`
}

type TxOutResp struct {
	TxIdHex       string `json:"txid"`
	Vout          int    `json:"vout"`
	Address       string `json:"address"`
	GenesisHex    string `json:"genesis"`
	Satoshi       int    `json:"satoshi"`
	ScriptTypeHex string `json:"script_type"`
	ScriptHex     string `json:"script"`
	Height        int    `json:"height"`
}

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
