package model

import (
	"encoding/binary"
	"encoding/json"

	scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
)

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
type TxoData struct {
	UTxid       []byte
	Vout        uint32
	BlockHeight uint32
	TxIdx       uint64
	Satoshi     uint64
	ScriptType  []byte
	Script      []byte
}

func (d *TxoData) Marshal(buf []byte) {
	binary.LittleEndian.PutUint32(buf, d.BlockHeight)  // 4
	binary.LittleEndian.PutUint64(buf[4:], d.TxIdx)    // 8
	binary.LittleEndian.PutUint64(buf[12:], d.Satoshi) // 8
	copy(buf[20:], d.Script)                           // n
}

func (d *TxoData) Unmarshal(buf []byte) {
	d.BlockHeight = binary.LittleEndian.Uint32(buf[:4]) // 4
	d.TxIdx = binary.LittleEndian.Uint64(buf[4:12])     // 8
	d.Satoshi = binary.LittleEndian.Uint64(buf[12:20])  // 8
	d.Script = buf[20:]
}

func NewTxoData(outpoint, res []byte) (txout *TxoData) {
	txout = &TxoData{}
	txout.Unmarshal(res)

	// 补充数据
	txout.UTxid = outpoint[:32]                            // 32
	txout.Vout = binary.LittleEndian.Uint32(outpoint[32:]) // 4
	txout.ScriptType = scriptDecoder.GetLockingScriptType(txout.Script)
	return
}
