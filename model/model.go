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
	PkScript    []byte
}

func (d *TxoData) Unmarshal(buf []byte) {
	if buf[3] == 0x00 {
		// not compress
		d.BlockHeight = binary.LittleEndian.Uint32(buf[:4]) // 4
		d.TxIdx = binary.LittleEndian.Uint64(buf[4:12])     // 8
		d.Satoshi = binary.LittleEndian.Uint64(buf[12:20])  // 8
		d.PkScript = buf[20:]
		return
	}

	buf[3] = 0x00
	d.BlockHeight = binary.LittleEndian.Uint32(buf[:4]) // 4

	offset := 4
	txidx, bytesRead := scriptDecoder.DeserializeVLQ(buf[offset:])
	if bytesRead >= len(buf[offset:]) {
		// errors.New("unexpected end of data after txidx")
		return
	}
	d.TxIdx = txidx

	offset += bytesRead
	compressedAmount, bytesRead := scriptDecoder.DeserializeVLQ(buf[offset:])
	if bytesRead >= len(buf[offset:]) {
		// errors.New("unexpected end of data after compressed amount")
		return
	}

	offset += bytesRead
	// Decode the compressed script size and ensure there are enough bytes
	// left in the slice for it.
	scriptSize := scriptDecoder.DecodeCompressedScriptSize(buf[offset:])
	if len(buf[offset:]) < scriptSize {
		// errors.New("unexpected end of data after script size")
		return
	}

	d.Satoshi = scriptDecoder.DecompressTxOutAmount(compressedAmount)
	d.PkScript = scriptDecoder.DecompressScript(buf[offset : offset+scriptSize])

}

func NewTxoData(outpoint, res []byte) (txout *TxoData) {
	txout = &TxoData{}
	txout.Unmarshal(res)

	// 补充数据
	txout.UTxid = outpoint[:32]                            // 32
	txout.Vout = binary.LittleEndian.Uint32(outpoint[32:]) // 4
	txout.ScriptType = scriptDecoder.GetLockingScriptType(txout.PkScript)
	return
}
