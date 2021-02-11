package blkparser

import (
	"encoding/binary"
)

type Tx struct {
	Hash     []byte // 32
	Size     uint32
	LockTime uint32
	Version  uint32
	TxInCnt  uint32
	TxOutCnt uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
}

type TxIn struct {
	InputHash []byte // 32
	InputVout uint32
	ScriptSig []byte
	Sequence  uint32
}

type TxOut struct {
	Value    uint64
	ValueRaw []byte
	Pkscript []byte
}

func NewTx(rawtx []byte) (tx *Tx, offset uint) {
	txLen := len(rawtx)
	if txLen < 4+1+32+4+1+1+1+8+1+1+4 {
		return nil, 0
	}

	tx = new(Tx)
	tx.Version = binary.LittleEndian.Uint32(rawtx[0:4])
	offset = 4

	txincnt, txincntsize := DecodeVariableLengthInteger(rawtx[offset:])
	offset += txincntsize

	tx.TxInCnt = uint32(txincnt)
	tx.TxIns = make([]*TxIn, txincnt)

	txoffset := uint(0)
	for i := range tx.TxIns {
		tx.TxIns[i], txoffset = NewTxIn(rawtx[offset:])
		// failed
		if txoffset == 0 {
			return nil, 0
		}
		offset += txoffset

		// invalid
		if offset >= uint(txLen) {
			return nil, 0
		}
	}

	txoutcnt, txoutcntsize := DecodeVariableLengthInteger(rawtx[offset:])
	offset += txoutcntsize

	tx.TxOutCnt = uint32(txoutcnt)
	tx.TxOuts = make([]*TxOut, txoutcnt)
	for i := range tx.TxOuts {
		tx.TxOuts[i], txoffset = NewTxOut(rawtx[offset:])
		// failed
		if txoffset == 0 {
			return nil, 0
		}

		offset += txoffset

		// invalid
		if offset >= uint(txLen) {
			return nil, 0
		}
	}

	// invalid
	if offset+4 != uint(txLen) {
		return nil, 0
	}

	tx.LockTime = binary.LittleEndian.Uint32(rawtx[offset : offset+4])
	offset += 4

	return
}

func NewTxIn(txinraw []byte) (txin *TxIn, offset uint) {
	inLen := len(txinraw)
	if inLen < 32+4+1+1+4 {
		return nil, 0
	}

	txin = new(TxIn)
	txin.InputHash = txinraw[0:32]
	txin.InputVout = binary.LittleEndian.Uint32(txinraw[32:36])
	offset = 36

	scriptsig, scriptsigsize := DecodeVariableLengthInteger(txinraw[offset:])
	offset += scriptsigsize

	// txin.ScriptSig = txinraw[offset : offset+scriptsig]
	offset += scriptsig

	// invalid
	if offset+4 > uint(inLen) {
		return nil, 0
	}

	txin.Sequence = binary.LittleEndian.Uint32(txinraw[offset : offset+4])
	offset += 4
	return
}

func NewTxOut(txoutraw []byte) (txout *TxOut, offset uint) {
	outLen := len(txoutraw)
	if outLen < 8+1+1 {
		return nil, 0
	}

	txout = new(TxOut)

	txout.ValueRaw = txoutraw[0:8]
	txout.Value = binary.LittleEndian.Uint64(txoutraw[0:8])
	offset = 8

	pkscript, pkscriptsize := DecodeVariableLengthInteger(txoutraw[offset:])
	offset += pkscriptsize

	// invalid
	if offset+pkscript > uint(outLen) {
		return nil, 0
	}

	txout.Pkscript = txoutraw[offset : offset+pkscript]
	offset += pkscript
	return
}
